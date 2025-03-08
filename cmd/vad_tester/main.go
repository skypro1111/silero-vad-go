package main

import (
	"encoding/binary"
	"flag"
	"fmt"
	"log/slog"
	"math"
	"os"
	"time"

	"github.com/streamer45/silero-vad-go/speech"
)

func main() {
	// Command line parameters
	modelPath := flag.String("model", "testfiles/silero_vad.onnx", "Path to Silero VAD model")
	audioPath := flag.String("audio", "testfiles/samples.pcm", "Path to PCM audio file")
	sampleRate := flag.Int("sr", 16000, "Sample rate (8000 or 16000)")
	threshold := flag.Float64("threshold", 0.5, "Speech detection probability threshold")
	negThreshold := flag.Float64("neg-threshold", 0.0, "Silence detection probability threshold (0 = auto)")
	minSilence := flag.Int("min-silence", 500, "Minimum silence duration (ms)")
	minSpeech := flag.Int("min-speech", 250, "Minimum speech duration (ms)")
	speechPad := flag.Int("speech-pad", 30, "Speech segments padding (ms)")
	verbose := flag.Bool("verbose", false, "Verbose output")
	flag.Parse()

	// Configure logging
	logLevel := slog.LevelInfo
	if *verbose {
		logLevel = slog.LevelDebug
	}
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: logLevel,
	}))
	slog.SetDefault(logger)

	// Load audio file
	slog.Info("Loading audio file", "path", *audioPath)
	samples, err := readPCMFile(*audioPath)
	if err != nil {
		slog.Error("Failed to load audio file", "error", err)
		os.Exit(1)
	}
	slog.Info("Audio file loaded", "samples", len(samples), "duration", fmt.Sprintf("%.2f sec", float64(len(samples))/float64(*sampleRate)))

	// Create detector configuration
	cfg := speech.DetectorConfig{
		ModelPath:            *modelPath,
		SampleRate:           *sampleRate,
		Threshold:            float32(*threshold),
		NegativeThreshold:    float32(*negThreshold),
		MinSilenceDurationMs: *minSilence,
		MinSpeechDurationMs:  *minSpeech,
		SpeechPadMs:          *speechPad,
		LogLevel:             speech.LogLevelError,
	}

	// Create detector
	slog.Info("Creating speech detector")
	startTime := time.Now()
	detector, err := speech.NewDetector(cfg)
	if err != nil {
		slog.Error("Failed to create detector", "error", err)
		os.Exit(1)
	}
	defer detector.Destroy()
	slog.Info("Detector created", "elapsed", time.Since(startTime))

	// Detect speech
	slog.Info("Starting speech detection", 
		"threshold", cfg.Threshold, 
		"negThreshold", cfg.NegativeThreshold,
		"minSilence", cfg.MinSilenceDurationMs,
		"minSpeech", cfg.MinSpeechDurationMs,
		"speechPad", cfg.SpeechPadMs)
	
	startTime = time.Now()
	segments, err := detector.Detect(samples)
	if err != nil {
		slog.Error("Speech detection failed", "error", err)
		os.Exit(1)
	}
	
	// Output results
	duration := time.Since(startTime)
	slog.Info("Speech detection completed", 
		"segments", len(segments), 
		"elapsed", duration,
		"rtf", duration.Seconds()/(float64(len(samples))/float64(*sampleRate)))

	fmt.Println("\nDetected speech segments:")
	fmt.Println("------------------------")
	totalSpeechDuration := 0.0
	for i, segment := range segments {
		segmentDuration := segment.SpeechEndAt - segment.SpeechStartAt
		if segment.SpeechEndAt > 0 {
			totalSpeechDuration += segmentDuration
			fmt.Printf("%d. %.2f - %.2f (%.2f sec)\n", i+1, segment.SpeechStartAt, segment.SpeechEndAt, segmentDuration)
		} else {
			fmt.Printf("%d. %.2f - [unfinished segment]\n", i+1, segment.SpeechStartAt)
		}
	}
	
	audioDuration := float64(len(samples)) / float64(*sampleRate)
	fmt.Printf("\nTotal audio duration: %.2f sec\n", audioDuration)
	fmt.Printf("Total speech duration: %.2f sec (%.1f%%)\n", 
		totalSpeechDuration, 
		(totalSpeechDuration/audioDuration)*100)
}

// Read PCM file with float32 samples
func readPCMFile(path string) ([]float32, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	samples := make([]float32, 0, len(data)/4)
	for i := 0; i < len(data); i += 4 {
		if i+4 <= len(data) {
			samples = append(samples, math.Float32frombits(binary.LittleEndian.Uint32(data[i:i+4])))
		}
	}
	return samples, nil
} 