package output

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
)

type Mode string

const (
	ModeJSON   Mode = "json"
	ModePretty Mode = "pretty"
)

type Renderer struct {
	Stdout io.Writer
	Stderr io.Writer
}

type envelope struct {
	OK    bool            `json:"ok"`
	Data  json.RawMessage `json:"data"`
	Error json.RawMessage `json:"error"`
	Meta  map[string]any  `json:"meta"`
}

func (r Renderer) WriteCapabilityResponse(body []byte, mode Mode, debug bool, format Format) error {
	var env envelope
	if err := json.Unmarshal(body, &env); err != nil {
		return err
	}

	if debug {
		debugPayload := map[string]any{
			"ok": env.OK,
		}
		if len(env.Data) > 0 && string(env.Data) != "null" {
			debugPayload["data"] = json.RawMessage(env.Data)
		}
		if len(env.Error) > 0 && string(env.Error) != "null" {
			debugPayload["error"] = json.RawMessage(env.Error)
		}
		if env.Meta != nil {
			delete(env.Meta, "source")
			debugPayload["meta"] = env.Meta
		}
		sanitized, err := json.Marshal(debugPayload)
		if err != nil {
			return err
		}
		return r.writeJSON(sanitized, mode)
	}

	if len(env.Error) > 0 && string(env.Error) != "null" {
		payload, err := json.Marshal(map[string]json.RawMessage{"error": env.Error})
		if err != nil {
			return err
		}
		return r.writeJSON(payload, mode)
	}
	if len(env.Data) == 0 {
		return r.writeJSON([]byte("null"), mode)
	}

	capability, _ := env.Meta["capability"].(string)
	projected, err := projectCapability(capability, env.Data, format)
	if err != nil {
		return err
	}
	return r.writeJSON(projected, mode)
}

func (r Renderer) writeJSON(body []byte, mode Mode) error {
	var out bytes.Buffer
	if mode == ModePretty {
		if err := json.Indent(&out, body, "", "  "); err != nil {
			return err
		}
		out.WriteByte('\n')
	} else {
		out.Write(body)
		if len(body) == 0 || body[len(body)-1] != '\n' {
			out.WriteByte('\n')
		}
	}
	_, err := r.Stdout.Write(out.Bytes())
	return err
}

func (r Renderer) Diagnostic(format string, args ...any) {
	if r.Stderr == nil {
		return
	}
	fmt.Fprintf(r.Stderr, format+"\n", args...)
}

func (r Renderer) WriteLocalError(message string, code string) error {
	body, err := json.Marshal(map[string]any{
		"error": map[string]any{
			"type":    "internal",
			"code":    code,
			"message": message,
		},
	})
	if err != nil {
		return err
	}
	return r.writeJSON(body, ModeJSON)
}
