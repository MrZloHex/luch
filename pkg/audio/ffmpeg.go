package audio

import (
	"errors"

	"github.com/3d0c/gmf"
)

type PCM struct {
	Samples []float32
	Rate    int
	Ch      int
}

func DecodeTo16kMono(path string) (PCM, error) {
	ictx, err := gmf.NewInputCtx(path)
	if err != nil {
		return PCM{}, err
	}
	defer ictx.CloseInputAndRelease()

	stream, err := ictx.GetBestStream(gmf.AVMEDIA_TYPE_AUDIO)
	if err != nil {
		return PCM{}, err
	}

	decCtx := stream.CodecCtx()
	if decCtx == nil {
		return PCM{}, errors.New("nil codec ctx")
	}

	// Open decoder
	codec, err := gmf.FindDecoder(decCtx.GetCodecId())
	if err != nil {
		return PCM{}, err
	}
	if err := decCtx.Open(codec, nil); err != nil {
		return PCM{}, err
	}
	defer decCtx.Free()

	// Setup resampler: target = s16, 16k, mono
	targetFmt := gmf.AV_SAMPLE_FMT_S16
	targetRate := 16000
	targetCh := 1
	swr, err := gmf.NewSwrCtx(
		decCtx.Channels(), decCtx.SampleFmt(), decCtx.SampleRate(),
		targetCh, targetFmt, targetRate,
	)
	if err != nil {
		return PCM{}, err
	}
	defer swr.Free()

	// FIFO to gather s16 samples
	fifo := gmf.NewAVAudioFifo(targetFmt, targetCh, 1024)
	defer fifo.Free()

	pkt := gmf.NewPacket()
	defer pkt.Free()

	for ictx.GetNextPacket(pkt) == 0 {
		if pkt.StreamIndex() != stream.Index() {
			pkt.Free()
			pkt = gmf.NewPacket()
			continue
		}
		frames, err := decCtx.Decode(pkt)
		pkt.Free()
		pkt = gmf.NewPacket()
		if err != nil && !gmf.IsAgain(err) {
			return PCM{}, err
		}
		for _, f := range frames {
			// Resample to target
			rf, err := swr.Convert(f)
			f.Free()
			if err != nil {
				return PCM{}, err
			}
			if rf == nil {
				continue
			}
			if err := fifo.WriteFrame(rf); err != nil {
				rf.Free()
				return PCM{}, err
			}
			rf.Free()
		}
	}

	// Drain decoder
	frames, _ := decCtx.Decode(nil)
	for _, f := range frames {
		rf, err := swr.Convert(f)
		f.Free()
		if err != nil {
			return PCM{}, err
		}
		if rf != nil {
			_ = fifo.WriteFrame(rf)
			rf.Free()
		}
	}

	// Pull all s16 samples
	var out []float32
	for fifo.Size() > 0 {
		rf, err := fifo.Read(1024)
		if err != nil {
			break
		}
		// rf.Data(0) is []byte containing interleaved s16 mono
		data := rf.Data(0)
		for i := 0; i+1 < len(data); i += 2 {
			// little-endian int16 -> float32 [-1,1]
			v := int16(uint16(data[i]) | uint16(data[i+1])<<8)
			out = append(out, float32(v)/32768.0)
		}
		rf.Free()
	}

	return PCM{Samples: out, Rate: targetRate, Ch: targetCh}, nil
}
