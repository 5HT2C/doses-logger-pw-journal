package types

import (
	"encoding/json"
	"fmt"
	"os"
)

const (
	ColorCfg = "config/colors.json"
)

var (
	DefaultColors = AdaptiveColorConfig{}
)

// AdaptiveColorConfig stored as [[Color]][AdaptiveColor].
//
// Contains all loaded [Journal]-defined colors from [ColorCfg].
type AdaptiveColorConfig map[Color]AdaptiveColor

// AdaptiveColor is a [Journal]-defined RGB color value which is used to tint the background of substances, as well as the timeline logs.
//   - Dark is the respective [ColorValue] for the app's dark theme
//   - Light is the respective [ColorValue] for the app's light theme
//   - Preferred is used in [FinishIngestionScreenViewModel.kt] to choose a random color when creating a new [JournalCompanion] for a [Substance] which does not already have its own [JournalCompanion].
//
// [FinishIngestionScreenViewModel.kt]: https://github.com/pwarchive/psychonautwiki-journal-android/blob/73f013752cea2f05558c1ed091cdccd3dfcde62b/app/src/main/java/com/isaakhanimann/journal/ui/tabs/journal/addingestion/time/FinishIngestionScreenViewModel.kt#L197
type AdaptiveColor struct {
	Dark      ColorValue `json:"dark"`
	Light     ColorValue `json:"light"`
	Preferred bool       `json:"preferred"`
}

// Color is a [Journal]-defined string enum name for each [ColorValue]
type Color string

// ColorValue is the stored RGB value of a [Color]
type ColorValue struct {
	Red   int64 `json:"r"`
	Green int64 `json:"g"`
	Blue  int64 `json:"b"`
}

// InitColors is used to read the [ColorCfg] JSON file which is generated using [AdaptiveColor.kts], which is patched from [AdaptiveColor.kt].
//
// For more information regarding the generation process of the const color types, refer to here:
// [pkg/github.com/5HT2C/doses-logger-pw-journal/gen.Render]
//
// [AdaptiveColor.kts]: https://github.com/5HT2C/doses-logger-pw-journal/tree/master/gen
// [AdaptiveColor.kt]: https://github.com/pwarchive/psychonautwiki-journal-android/blob/73f013752cea2f05558c1ed091cdccd3dfcde62b/app/src/main/java/com/isaakhanimann/journal/data/room/experiences/entities/AdaptiveColor.kt#L587
func InitColors() error {
	if b, err := os.ReadFile(ColorCfg); err != nil {
		return err
	} else {
		return json.Unmarshal(b, &DefaultColors)
	}
}

// ColorValue returns a formatted [R, G, B] string.
//
// Example:
//
//	fmt.Printf("%s\n", *types.ColorCyan.Adapt(true))
//	// [100, 210, 255]
func (c ColorValue) String() string {
	return fmt.Sprintf("[%03d, %03d, %03d]", c.Red, c.Green, c.Blue)
}

// AdaptiveColor returns a formatted [R, G, B] string for each respective dark and light [ColorValue].
//
// Example:
//
//	fmt.Println(types.ColorCyan.Value().String())
//	// D: [100, 210, 255] | L: [050, 173, 230] | P: true
func (c AdaptiveColor) String() string {
	return fmt.Sprintf(
		"D: %s | L: %s | P: %v",
		c.Dark.String(), c.Light.String(), c.Preferred,
	)
}

// String returns a formatted [AdaptiveColor] string if it exists for the given [Color].
// Otherwise, it returns a formatted error message.
func (c Color) String() string {
	if v := c.Value(); v != nil {
		return v.String()
	}

	return "nil (nonexistent color value!)"
}

// Adapt provides returns the RGB [ColorValue] for a given theme,
// as a convenience method when using the generated const types
//
// Preferably use this instead of [Color.Adapt].
func (c AdaptiveColor) Adapt(darkTheme bool) ColorValue {
	if darkTheme { // TODO: This can probably be entirely reworked
		return c.Dark
	} else {
		return c.Light
	}
}

// Value returns the internal [AdaptiveColor] value for a declared color.
func (c Color) Value() *AdaptiveColor {
	if v, ok := DefaultColors[c]; ok {
		return &v
	}

	return nil
}

// Adapt provides returns the RGB [ColorValue] for a given theme,
// as a convenience method when using the generated const types.
//
// Do not use this unless you must, as it does not guarantee that a valid color value will be returned.
// If the color name is misspelled or does not exist in [DefaultColors] it will return nil.
//
// Preferably use [AdaptiveColor.Adapt] instead.
func (c Color) Adapt(darkTheme bool) *ColorValue {
	if v := c.Value(); v != nil {
		vc := v.Adapt(darkTheme)
		return &vc
	}

	return nil
}
