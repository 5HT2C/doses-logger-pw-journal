# Building

Below are instructions to generate the `config/colors.json` file using PsychonautWiki Journal's original source code.

1. Download `AdaptiveColor.kt` from a mirror[^1] such as [this one](https://github.com/pwarchive/psychonautwiki-journal-android/blob/73f013752cea2f05558c1ed091cdccd3dfcde62b/app/src/main/java/com/isaakhanimann/journal/data/room/experiences/entities/AdaptiveColor.kt) (you can use "Raw" on GitHub to get a `$link` for the commands below)
2. Apply the patch file (ensure it is inside `gen/` along with `AdaptiveColor.kt`)
3. Rename `AdaptiveColor.kt` to `AdaptiveColor.kts`
4. Install and run `kotlin`[^2] in order to generate the JSON config

```shell
cd doses-logger-pw-journal/gen/

link='https://raw.githubusercontent.com/REPLACEME/REPLACEME/refs/heads/main/app/src/main/java/com/isaakhanimann/journal/data/room/experiences/entities/AdaptiveColor.kt'
curl -oL AdaptiveColor.kt "$link"
git apply AdaptiveColor.kt

mv AdaptiveColor.kt AdaptiveColor.kts
kotlin --script AdaptiveColor.kts > ../config/colors.json
```

[^1]: Unfortunately, the official `psychonautwiki-journal-android` sources have since been made private (as of version ≥`12.x`).
[^2]: [Kotlin command-line compiler | Kotlin Documentation](https://kotlinlang.org/docs/command-line.html#run-scripts)
