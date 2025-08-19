package audio

import (
    "fmt"

    "github.com/ggerganov/whisper.cpp/bindings/go/pkg/whisper"
)

type Engine struct
{
    ctx *whisper.Context
}

func NewEngine(modelPath string) (*Engine, error)
{
    ctx, err := whisper.NewContext(modelPath)
    if err != nil { return nil, err }
    return &Engine{ctx: ctx}, nil
}

func (e *Engine) Close()
{
    if e.ctx != nil { e.ctx.Close() }
}

type Params struct
{
    Language   string // "auto" or "en"/"ru"/...
    Translate  bool   // true: force translate to English
    NoTimestamps bool
}

func (e *Engine) TranscribeFloat32(samples []float32, p Params) (string, error)
{
    wp := whisper.NewFullParams(whisper.SAMPLING_GREEDY)
    if p.Language == "" { p.Language = "auto" }
    if p.Language != "auto" {
        if err := wp.SetLanguage(p.Language); err != nil { return "", err }
    } else {
        _ = wp.SetDetectLanguage(true)
    }
    _ = wp.SetTranslate(p.Translate)
    _ = wp.SetNoTimestamps(p.NoTimestamps)

    if err := e.ctx.Process(samples, wp); err != nil {
        return "", err
    }

    // Collect segments
    n := e.ctx.NumSegments()
    if n == 0 { return "", fmt.Errorf("no segments") }

    out := make([]byte, 0, 1024)
    for i := 0; i < n; i++ {
        seg := e.ctx.GetSegment(i)
        out = append(out, seg.Text()...)
        if i+1 < n { out = append(out, ' ') }
    }
    return string(out), nil
}

