<h1 align="center">
  <br>
  silero-vad-go
  <br>
</h1>
<h4 align="center">A simple Golang (CGO + ONNX Runtime) speech detector powered by Silero VAD</h4>
<p align="center">
  <a href="https://pkg.go.dev/github.com/skypro1111/silero-vad-go"><img src="https://pkg.go.dev/badge/github.com/skypro1111/silero-vad-go.svg" alt="Go Reference"></a>
  <a href="LICENSE"><img src="https://img.shields.io/badge/License-MIT-yellow.svg" alt="License: MIT"></a>
</p>
<br>

# Silero VAD Go

Go implementation of the [Silero VAD](https://github.com/snakers4/silero-vad) (Voice Activity Detection) using ONNX Runtime.

## Requirements

- [Golang](https://go.dev/doc/install) >= v1.21
- A C compiler (e.g. GCC)
- ONNX Runtime (v1.18.1)
- A [Silero VAD](https://github.com/snakers4/silero-vad) model (v5)

### Installing ONNX Runtime

#### macOS
```bash
brew install onnxruntime
```

#### Linux
Download and install from [ONNX Runtime releases](https://github.com/microsoft/onnxruntime/releases)

For development, you need to export the following environment variables:

#### Linux
```sh
LD_RUN_PATH="/usr/local/lib/onnxruntime-linux-x64-1.18.1/lib"
LIBRARY_PATH="/usr/local/lib/onnxruntime-linux-x64-1.18.1/lib"
C_INCLUDE_PATH="/usr/local/include/onnxruntime-linux-x64-1.18.1/include"
```

#### Darwin (MacOS)
```sh
LIBRARY_PATH="/usr/local/lib/onnxruntime-linux-x64-1.18.1/lib"
C_INCLUDE_PATH="/usr/local/include/onnxruntime-linux-x64-1.18.1/include"
sudo update_dyld_shared_cache
```

## Installation

```bash
go get github.com/skypro1111/silero-vad-go
```

## Quick Start

1. Download the Silero VAD model in ONNX format:
```bash
wget https://github.com/snakers4/silero-vad/raw/master/files/silero_vad.onnx
```

2. Create a new detector:
```go
package main

import "github.com/skypro1111/silero-vad-go/speech"

func main() {
    // Create detector configuration
    config := speech.DetectorConfig{
        ModelPath:          "silero_vad.onnx",
        SampleRate:         16000,
        Threshold:          0.5,
        NegativeThreshold:  0.0,  // Auto: threshold - 0.15
        MinSilenceDuration: 500,  // milliseconds
        MinSpeechDuration:  250,  // milliseconds
        SpeechPadMs:       30,    // milliseconds
    }

    // Create new detector
    detector, err := speech.NewDetector(config)
    if err != nil {
        panic(err)
    }
    defer detector.Close()

    // Your audio processing code here...
}
```

## Usage Examples

### Basic Audio Processing

```go
// Read PCM audio data (16-bit integers)
audioData := make([]int16, sampleRate*duration) // e.g., 16000 samples for 1 second
// ... fill audioData with your audio samples ...

// Process audio in chunks
chunkSize := 512
for i := 0; i < len(audioData); i += chunkSize {
    end := i + chunkSize
    if end > len(audioData) {
        end = len(audioData)
    }
    chunk := audioData[i:end]
    
    // Process chunk
    isSpeech, err := detector.Process(chunk)
    if err != nil {
        panic(err)
    }
    
    if isSpeech {
        fmt.Println("Speech detected!")
    }
}

// Get speech segments
segments := detector.GetSegments()
for _, seg := range segments {
    fmt.Printf("Speech segment: %.2f - %.2f (%.2f sec)\n", 
        seg.Start, seg.End, seg.End-seg.Start)
}
```

### Real-time Audio Processing

```go
// Example with PortAudio
import "github.com/gordonklaus/portaudio"

func main() {
    config := speech.DetectorConfig{
        ModelPath:          "silero_vad.onnx",
        SampleRate:         16000,
        Threshold:          0.5,
    }
    
    detector, err := speech.NewDetector(config)
    if err != nil {
        panic(err)
    }
    defer detector.Close()

    // Initialize PortAudio
    if err := portaudio.Initialize(); err != nil {
        panic(err)
    }
    defer portaudio.Terminate()

    // Open default input stream
    inputBuffer := make([]int16, 512)
    stream, err := portaudio.OpenDefaultStream(1, 0, float64(config.SampleRate), 
        len(inputBuffer), inputBuffer)
    if err != nil {
        panic(err)
    }
    defer stream.Close()

    if err := stream.Start(); err != nil {
        panic(err)
    }

    // Process audio in real-time
    for {
        if err := stream.Read(); err != nil {
            panic(err)
        }

        isSpeech, err := detector.Process(inputBuffer)
        if err != nil {
            panic(err)
        }

        if isSpeech {
            fmt.Println("Speech detected!")
        }
    }
}
```

### Parameter Tuning Examples

1. More sensitive detection (for quiet speech):
```go
config := speech.DetectorConfig{
    ModelPath:          "silero_vad.onnx",
    SampleRate:         16000,
    Threshold:          0.3,
    NegativeThreshold:  0.1,
    MinSpeechDuration:  150,
}
```

2. Conservative detection (reduce false positives):
```go
config := speech.DetectorConfig{
    ModelPath:          "silero_vad.onnx",
    SampleRate:         16000,
    Threshold:          0.7,
    MinSilenceDuration: 1000,
    MinSpeechDuration:  500,
}
```

3. Quick response (for short commands):
```go
config := speech.DetectorConfig{
    ModelPath:          "silero_vad.onnx",
    SampleRate:         16000,
    Threshold:          0.5,
    NegativeThreshold:  0.4,
    MinSilenceDuration: 200,
    MinSpeechDuration:  200,
    SpeechPadMs:       10,
}
```

4. Noisy environment:
```go
config := speech.DetectorConfig{
    ModelPath:          "silero_vad.onnx",
    SampleRate:         16000,
    Threshold:          0.6,
    NegativeThreshold:  0.45,
    MinSilenceDuration: 700,
    MinSpeechDuration:  300,
    SpeechPadMs:       50,
}
```

## Testing

The repository includes a test utility (`cmd/vad_tester`) for experimenting with different parameters. See [VAD Tester README](cmd/vad_tester/README.md) for detailed usage instructions.

Quick test example:
```bash
cd cmd/vad_tester
go build
./vad_tester -model ../../testfiles/silero_vad.onnx -audio ../../testfiles/samples.pcm
```

## Parameters

### Speech Detection Threshold
- Range: (0, 1)
- Default: 0.5
- Higher values: more conservative, less false positives
- Lower values: more sensitive, might catch quieter speech

### Silence Detection Threshold
- Range: [0, 1)
- Default: 0.0 (auto: threshold - 0.15)
- Higher values: faster switching to silence
- Lower values: more persistent speech detection

### Minimum Silence Duration
- Unit: milliseconds
- Default: 500
- Higher values: fewer splits, longer segments
- Lower values: more splits, shorter segments

### Minimum Speech Duration
- Unit: milliseconds
- Default: 250
- Higher values: filter out short noises
- Lower values: catch brief utterances

### Speech Padding
- Unit: milliseconds
- Default: 30
- Higher values: less aggressive cutting
- Lower values: tighter segments

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

