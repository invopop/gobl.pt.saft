package saft_test

import (
	"bytes"
	"strings"
	"testing"

	saft "github.com/invopop/gobl.pt.saft"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/text/encoding/charmap"
)

type testStruct struct {
	Name        string
	Value       string
	Description string
}

func TestMarshalXML(t *testing.T) {
	t.Run("Simple struct without root", func(t *testing.T) {
		xmlData, err := saft.MarshalXML(testStruct{Name: "Test", Value: "123"}, false)
		require.NoError(t, err)
		require.NotEmpty(t, xmlData)

		// Verify it's valid Windows-1252
		decoded, err := decodeWindows1252(xmlData)
		require.NoError(t, err)
		assert.Contains(t, decoded, "<Name>Test</Name>")
		assert.Contains(t, decoded, "<Value>123</Value>")
	})

	t.Run("With root header", func(t *testing.T) {
		xmlData, err := saft.MarshalXML(testStruct{Name: "Test"}, true)
		require.NoError(t, err)
		require.NotEmpty(t, xmlData)

		// Should start with XML declaration
		decoded, err := decodeWindows1252(xmlData)
		require.NoError(t, err)
		assert.True(t, strings.HasPrefix(decoded, `<?xml version="1.0" encoding="Windows-1252"?>`))
	})

	t.Run("NFC normalization preserves accents", func(t *testing.T) {
		nfdText := "P\u0065\u0302ssego" // "Pêssego" using NFD (decomposed)
		xmlData, err := saft.MarshalXML(testStruct{Name: nfdText}, false)
		require.NoError(t, err)
		require.NotEmpty(t, xmlData)

		// Verify accent is preserved
		decoded, err := decodeWindows1252(xmlData)
		require.NoError(t, err)
		assert.Contains(t, decoded, "Pêssego")
	})

	t.Run("Unsupported characters do not cause errors", func(t *testing.T) {
		unsupportedText := "HERDADE DA MESQUITA，ESTRADA" // Fullwidth comma
		xmlData, err := saft.MarshalXML(testStruct{Name: unsupportedText}, false)
		require.NoError(t, err)
		require.NotEmpty(t, xmlData)

		// Verify output is valid Windows-1252
		decoded, err := decodeWindows1252(xmlData)
		require.NoError(t, err)
		require.NotEmpty(t, decoded)
		assert.Contains(t, decoded, "HERDADE DA MESQUITA&#65292;ESTRADA")
	})
}

func TestUnmarshalXML(t *testing.T) {
	t.Run("Decode Windows-1252 XML", func(t *testing.T) {
		// Create XML in Windows-1252 with precomposed ê (0xEA)
		xmlStr := "<testStruct><Name>Pêssego</Name></testStruct>"
		xmlBytes := encodeToWindows1252(t, xmlStr)

		var result testStruct
		err := saft.UnmarshalXML(xmlBytes, &result)
		require.NoError(t, err)
		assert.Equal(t, "Pêssego", result.Name)
	})

	t.Run("Round trip marshal and unmarshal", func(t *testing.T) {
		original := testStruct{
			Name:        "Pêssego Fuze Tea",
			Description: "Bebida de chá com pêssego",
		}

		// Marshal
		xmlData, err := saft.MarshalXML(original, false)
		require.NoError(t, err)

		// Unmarshal
		var decoded testStruct
		err = saft.UnmarshalXML(xmlData, &decoded)
		require.NoError(t, err)

		// Should match original
		assert.Equal(t, original.Name, decoded.Name)
		assert.Equal(t, original.Description, decoded.Description)
	})
}

func TestEncodeXML(t *testing.T) {
	t.Run("Write to buffer without root", func(t *testing.T) {
		var buf bytes.Buffer

		err := saft.EncodeXML(testStruct{Name: "Test"}, &buf, false)
		require.NoError(t, err)

		decoded, err := decodeWindows1252(buf.Bytes())
		require.NoError(t, err)
		assert.Contains(t, decoded, "<Name>Test</Name>")
	})

	t.Run("Write to buffer with root", func(t *testing.T) {
		var buf bytes.Buffer

		err := saft.EncodeXML(testStruct{Name: "Test"}, &buf, true)
		require.NoError(t, err)

		decoded, err := decodeWindows1252(buf.Bytes())
		require.NoError(t, err)
		assert.True(t, strings.HasPrefix(decoded, `<?xml version="1.0" encoding="Windows-1252"?>`))
		assert.Contains(t, decoded, "<Name>Test</Name>")
	})

	t.Run("NFD normalization in encoder", func(t *testing.T) {
		var buf bytes.Buffer

		err := saft.EncodeXML(testStruct{Name: "P\u0065\u0302ssego"}, &buf, false)
		require.NoError(t, err)

		decoded, err := decodeWindows1252(buf.Bytes())
		require.NoError(t, err)
		assert.Contains(t, decoded, "Pêssego")
	})
}

// Helper function to decode Windows-1252 bytes to UTF-8 string
func decodeWindows1252(data []byte) (string, error) {
	decoder := charmap.Windows1252.NewDecoder()
	decoded, err := decoder.Bytes(data)
	if err != nil {
		return "", err
	}
	return string(decoded), nil
}

// Helper function to encode UTF-8 string to Windows-1252 bytes
func encodeToWindows1252(t *testing.T, s string) []byte {
	encoder := charmap.Windows1252.NewEncoder()
	encoded, err := encoder.String(s)
	require.NoError(t, err)
	return []byte(encoded)
}
