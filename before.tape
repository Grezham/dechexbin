# VHS documentation
#
# Output:
#   Output <path>.gif               Create a GIF output at the given <path>
#   Output <path>.mp4               Create an MP4 output at the given <path>
#   Output <path>.webm              Create a WebM output at the given <path>
#
# Require:
#   Require <string>                Ensure a program is on the $PATH to proceed
#
# Settings:
#   Set FontSize <number>           Set the font size of the terminal
#   Set FontFamily <string>         Set the font family of the terminal
#   Set Height <number>             Set the height of the terminal
#   Set Width <number>              Set the width of the terminal
#   Set LetterSpacing <float>       Set the font letter spacing (tracking)
#   Set LineHeight <float>          Set the font line height
#   Set LoopOffset <float>%         Set the starting frame offset for the GIF loop
#   Set Theme <json|string>         Set the theme of the terminal
#   Set Padding <number>            Set the padding of the terminal
#   Set Framerate <number>          Set the framerate of the recording
#   Set PlaybackSpeed <float>       Set the playback speed of the recording
#
# Sleep:
#   Sleep <time>                    Sleep for a set amount of <time> in seconds
#
# Type:
#   Type[@<time>] "<characters>"    Type <characters> into the terminal with a
#                                   <time> delay between each character
#
# Keys:
#   Backspace[@<time>] [number]     Press the Backspace key
#   Down[@<time>] [number]          Press the Down key
#   Enter[@<time>] [number]         Press the Enter key
#   Space[@<time>] [number]         Press the Space key
#   Tab[@<time>] [number]           Press the Tab key
#   Left[@<time>] [number]          Press the Left Arrow key
#   Right[@<time>] [number]         Press the Right Arrow key
#   Up[@<time>] [number]            Press the Up Arrow key
#   Down[@<time>] [number]          Press the Down Arrow key
#   PageUp[@<time>] [number]        Press the Page Up key
#   PageDown[@<time>] [number]      Press the Page Down key
#   Ctrl+<key>                      Press the Control key + <key> (e.g. Ctrl+C)
#
# Display:
#   Hide                            Hide the subsequent commands from the output
#   Show                            Show the subsequent commands in the output

Output before.gif

Require ./main

Set Shell "bash"
Set FontSize 32
Set Width 1400
Set Height 700
#Enter Program
Type@500ms "./main" Enter
Sleep 1s
Type@500ms "4"
Down 1
Type@500ms "5"
Sleep 1s
Down 1
Right@800ms 2
Sleep 1s
Down 1
Right@800ms 2
Left@800ms 1
Sleep 1s
Enter

#Question Set
Type@500ms "1"
Sleep 1s Enter
Sleep 2s
Type@500ms "200"
Sleep 1s Enter
Type@500ms "Let's just get to the end Review"
Enter@1000ms 2
Sleep 2s
Type "?"
Sleep 2s
Type "r"
Sleep 1s
Type@800ms "'r' Restarts the Quiz with randomize questions"
Sleep 800ms
Backspace@500ms 47
Sleep 500ms
Type@800ms "And 's' takes us back to the Settings view"
Sleep 800ms
Type "s"
Sleep 1s
Type "q"

Sleep 2s
