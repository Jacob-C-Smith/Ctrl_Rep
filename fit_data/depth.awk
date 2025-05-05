BEGIN {
    FS = "\\?"
    OFS = " ? "
}

{
    key = $1
    gsub(/^[ \t]+|[ \t]+$/, "", key)
    gsub(/^[ \t]+|[ \t]+$/, "", $2)

    # Count ':' in key
    numColons = gsub(/:/, ":", key)

    if (numColons <= depth) {
        keys[NR] = key
        values[NR] = $2
        lengths[NR] = length(key)
        if (length(key) > maxLen) {
            maxLen = length(key)
        }
    }
}

END {
    for (i = 1; i <= NR; i++) {
        if (i in keys) {
            printf "%-*s ? %s\n", maxLen, keys[i], values[i]
        }
    }
}
