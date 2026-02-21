package tracks

import (
	"errors"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/hraban/opus"
	"github.com/pion/rtp"
	webrtc "github.com/pion/webrtc/v4"
)

// MySineTrack implements a TrackLocal interface
type MySineTrack struct {
	streamId  string
	sampleIdx int
	samples   []int16
}

func NewMySineTrack() *MySineTrack {
	return &MySineTrack{
		streamId: uuid.New().String(),
	}
}

func getOpusCodecParams(ctx webrtc.TrackLocalContext) *webrtc.RTPCodecParameters {
	for _, codec := range ctx.CodecParameters() {
		if codec.RTPCodecCapability.MimeType == webrtc.MimeTypeOpus {
			return &codec
		}
	}
	return nil
}

func (t *MySineTrack) startStreaming(ctx webrtc.TrackLocalContext, selectedCodec webrtc.RTPCodecParameters) {
	// 1. Initialize the Encoder (48kHz, 2 Channels for Stereo)
	// For Mono, use 1 channel and adjust your sample logic.
	enc, err := opus.NewEncoder(48000, 2, opus.AppVoIP)
	if err != nil {
		return
	}

	ticker := time.NewTicker(20 * time.Millisecond)
	samplesPerPacket := 960
	var sequenceNumber uint16
	var timestamp uint32

	for {
		select {
		case <-ticker.C:
			// 2. Grab the 960 samples (Stereo = 1920 int16s)
			end := t.sampleIdx + (samplesPerPacket * 2)
			if end > len(t.samples) {
				t.sampleIdx = 0
				end = samplesPerPacket * 2
			}
			pcmChunk := t.samples[t.sampleIdx:end]

			// 3. Encode PCM to Opus
			data := make([]byte, 1000) // Buffer for compressed data
			n, err := enc.Encode(pcmChunk, data)
			if err != nil {
				continue
			}

			rtpHeader := rtp.Header{
				Version:        2,
				PayloadType:    uint8(selectedCodec.PayloadType),
				SequenceNumber: sequenceNumber,
				Timestamp:      timestamp,
				SSRC:           uint32(ctx.SSRC()),
				Marker:         false,
			}

			// 5. Send and Increment
			if _, err := ctx.WriteStream().WriteRTP(&rtpHeader, data[:n]); err == nil {
				sequenceNumber++
				timestamp += uint32(samplesPerPacket)
				t.sampleIdx += (samplesPerPacket * 2)
			}

		case <-ctx.Done():
			return
		}
	}
}

// Bind should implement the way how the media data flows from the Track to the PeerConnection
// This will be called internally after signaling is complete and the list of available
// codecs has been determined
func (track *MySineTrack) Bind(ctx webrtc.TrackLocalContext) (webrtc.RTPCodecParameters, error) {
	codecParam := getOpusCodecParams(ctx)
	if codecParam == nil {
		return webrtc.RTPCodecParameters{}, errors.New("no supported codec found, currently only opus is supported")
	}

	track.startStreaming(ctx)

	return *codecParam, nil
}

// Unbind should implement the teardown logic when the track is no longer needed. This happens
// because a track has been stopped.
func (track *MySineTrack) Unbind(ctx webrtc.TrackLocalContext) error {
	// todo: teardown and clean up
	log.Printf("[track] stream %s is tearing down", track.streamId)
	return nil
}

// ID is the unique identifier for this Track. This should be unique for the
// stream, but doesn't have to globally unique. A common example would be 'audio' or 'video'
// and StreamID would be 'desktop' or 'webcam'
func (track *MySineTrack) ID() string {
	return "audio"
}

// RID is the RTP Stream ID for this track.
func (track *MySineTrack) RID() string {
	return "my_sine_track_rid"
}

// StreamID is the group this track belongs too. This must be unique
func (track *MySineTrack) StreamID() string {
	return track.streamId
}

// Kind controls if this TrackLocal is audio or video
func (track *MySineTrack) Kind() webrtc.RTPCodecType {
	return webrtc.RTPCodecTypeAudio
}
