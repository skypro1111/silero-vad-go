#!/bin/bash

# Make sure the program is compiled
echo "Compiling program..."
go build -o vad_tester

# Base parameters
MODEL="../../testfiles/silero_vad.onnx"
AUDIO="../../testfiles/samples2.pcm"
SR=16000

# Function to run test with specific parameters
run_test() {
    local name=$1
    local threshold=$2
    local neg_threshold=$3
    local min_silence=$4
    local min_speech=$5
    local speech_pad=$6
    
    echo "============================================================"
    echo "Test: $name"
    echo "------------------------------------------------------------"
    echo "Parameters:"
    echo "  Threshold: $threshold"
    echo "  Negative Threshold: $neg_threshold"
    echo "  Min Silence Duration: $min_silence ms"
    echo "  Min Speech Duration: $min_speech ms"
    echo "  Speech Pad: $speech_pad ms"
    echo "------------------------------------------------------------"
    
    ./vad_tester \
        -model $MODEL \
        -audio $AUDIO \
        -sr $SR \
        -threshold $threshold \
        -neg-threshold $neg_threshold \
        -min-silence $min_silence \
        -min-speech $min_speech \
        -speech-pad $speech_pad
    
    echo "============================================================"
    echo ""
}

# Test 1: Base parameters
run_test "Base parameters" 0.5 0.0 500 250 30

# Test 2: Low speech threshold
run_test "Low speech threshold" 0.3 0.0 500 250 30

# Test 3: High speech threshold
run_test "High speech threshold" 0.7 0.0 500 250 30

# Test 4: Low negative threshold
run_test "Low negative threshold" 0.5 0.2 500 250 30

# Test 5: High negative threshold
run_test "High negative threshold" 0.5 0.4 500 250 30

# Test 6: Short minimum silence duration
run_test "Short minimum silence duration" 0.5 0.0 200 250 30

# Test 7: Long minimum silence duration
run_test "Long minimum silence duration" 0.5 0.0 1000 250 30

# Test 8: Short minimum speech duration
run_test "Short minimum speech duration" 0.5 0.0 500 100 30

# Test 9: Long minimum speech duration
run_test "Long minimum speech duration" 0.5 0.0 500 500 30

# Test 10: No padding for speech segments
run_test "No padding for speech segments" 0.5 0.0 500 250 0

# Test 11: Large padding for speech segments
run_test "Large padding for speech segments" 0.5 0.0 500 250 100

echo "All tests completed!" 