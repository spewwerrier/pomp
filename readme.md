# Pomp
Pomp is a pomodoro timer for waybar enjoyers.\\
There are 2 modes: working mode which I prefer to be long and break mode which is a short timer.
They can be toggled using signals. SIGUSR1 is for toggling state of pomp and SIGUSR2 is for
toggling the timer of the pomp.

You can choose your own working and rest mode by `pomp -w 7200 -r 1800`.\\

To create an executable
```sh
go build -o pomp main.go
```
Put pomp in your path. And don't forget to create waybar module for pomp in waybar config.

```json
"custom/pomo":{
	"exec": "/path/to/pomp",
	"return-type": "json",
	"format": "{icon} {text}",
	"on-click": "pkill -USR1 pomp",
	"on-click-right": "pkill -USR2 pomp",
	"tooltip-format": "A pomodoro timer",
	"format-icons": {
		"WORKING": "W",
		"RESTING": "R",
		"NON": "N",
	}
}
```
And consider styling it as well.
```css
#custom-pomo {
  color: white;
  background-color: black;
  border-radius: 10px;
}
```

Right click on module to start the timer on initial start. Otherwise toggle the state from
working to resting and vice versa. Left click pauses the timer.
