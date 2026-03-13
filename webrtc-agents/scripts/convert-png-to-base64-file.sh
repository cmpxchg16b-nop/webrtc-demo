#!/bin/bash

set -e

if [ -z "$1" ]; then
    exit 1
fi

dir_name=$(dirname $1)
base_name=$(basename $1)
output_basename="${base_name}.dataurl"

cd "$dir_name"


function determine_mime_from_suffix {
    local filename="$1"
    local extension="${filename##*.}"
    extension=$(echo "$extension" | tr '[:upper:]' '[:lower:]')

    local mime_type=""

    case "$extension" in
        png)
            mime_type="image/png"
            ;;
        svg)
            mime_type="image/svg+xml"
            ;;
        jpg|jpeg)
            mime_type="image/jpeg"
            ;;
        webp)
            mime_type="image/webp"
            ;;
        tif|tiff)
            mime_type="image/tiff"
            ;;
        *)
            mime_type=""
            ;;
    esac

    echo "$mime_type"
}

mime=$(determine_mime_from_suffix "$base_name")
if [ -z "$mime" ]; then
    echo "Unable to determine MIME type for $base_name"
    exit 1
fi

dataURLPrefix="data:$mime;base64,"

echo -n "$dataURLPrefix" > $output_basename
base64 -i $base_name >> $output_basename

echo Base64-encoded file has been written to $dir_name/$output_basename
