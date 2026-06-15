#!/usr/bin/env bash
#
# Generate a spectrogram image from an audio file. Either open it directly or
# output it to a file. Used to test whether audio files are truly lossless.

spectrogram() {
    local audio_file=""
    local output=""
    local open_after=1

    if [[ $# -eq 0 ]]; then
        echo "Usage: spectrogram [-o|--output [file]] <audio-file>"
        return 1
    fi

    # Audio file is treated as last argument
    audio_file="${@: -1}"

    if [[ ! -f "$audio_file" ]]; then
        echo "spectrogram: '$audio_file': not a file" >&2
        return 1
    fi

    while [[ $# -gt 1 ]]; do
        case "$1" in
            -o|--output)
                open_after=0
                # If an argument sits between this flag and the trailing input
                # file, use it as the output path ("-o file.png"). Otherwise
                # "-o" was given on its own and we fall back to a default name.
                if [[ $# -gt 2 ]]; then
                    output="$2"
                    shift
                else
                    output="spec.png"
                fi
                ;;
        esac
        shift
    done

    # Only force-overwrite (-y) the throwaway temp file we create ourselves. For
    # a user-supplied path, warn and let ffmpeg's own prompt handle it.
    local overwrite=()
    if [[ "$open_after" -eq 1 ]]; then
        output="$(mktemp --suffix=.png)"
        overwrite=(-y)
    fi

    ffmpeg -hide_banner -loglevel error \
        "${overwrite[@]}" \
        -i "$audio_file" \
        -lavfi "showspectrumpic=s=1920x1080:legend=1" \
        "$output" || return 1

    if [[ "$open_after" -eq 1 ]]; then
        xdg-open "$output"
    fi
}

spectrogram "$@"
