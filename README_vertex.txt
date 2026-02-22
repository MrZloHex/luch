███╗   ███╗ ██████╗ ███╗   ██╗ ██████╗ ██╗     ██╗████████╗██╗  ██╗
████╗ ████║██╔═══██╗████╗  ██║██╔═══██╗██║     ██║╚══██╔══╝██║  ██║
██╔████╔██║██║   ██║██╔██╗ ██║██║   ██║██║     ██║   ██║   ███████║
██║╚██╔╝██║██║   ██║██║╚██╗██║██║   ██║██║     ██║   ██║   ██╔══██║
██║ ╚═╝ ██║╚██████╔╝██║ ╚████║╚██████╔╝███████╗██║   ██║   ██║  ██║
╚═╝     ╚═╝ ╚═════╝ ╚═╝  ╚═══╝ ╚═════╝ ╚══════╝╚═╝   ╚═╝   ╚═╝  ╚═╝


  ░▒▓█ _vertex_ █▓▒░  
  The beating heart of **СоюзПечаль**-a night lamp that does more than shine.  

  ───────────────────────────────────────────────────────────────  
  ▓ OVERVIEW  
  **_vertex_** is a firmware running on an **ATmega328P**, controlling an  
  LED strip and a buzzer inside **СоюзПечаль**.  
  It communicates with an _obelisk_ via USART,  
  allowing dynamic lighting modes and interactions.  
  Simple. Flexible. Cool as hell.  

  ───────────────────────────────────────────────────────────────  
  ▓ HARDWARE  
  ▪ **MCU**: ATmega328P (Currently Arduino Nano, later standalone)  
  ▪ **OUTPUTS**: LED Strip + Buzzer  
  ▪ **INTERFACE**: USART for external command control  
  ▪ **POWER**: USB / External 5V  

  ───────────────────────────────────────────────────────────────  
  ▓ FEATURES  
  ▪ Custom lighting modes controlled via USART  
  ▪ Reactive buzzer for audio cues  
  ▪ Low-power idle mode  
  ▪ Designed for easy flashing & standalone ATmega328P use  

  ───────────────────────────────────────────────────────────────  
  ▓ BUILD & FLASH  
  Compile with **AVR-GCC**
  Flash using **AVRDUDE**

  ```sh  
  make &&
  make PORT=/dev/ttyUSB0 flash
  ```

  ───────────────────────────────────────────────────────────────  
  ▓ PROTOCOL  
  Packet format:  <TO>:<PAYLOAD>:<FROM>\n  
  Payload format: VERB:NOUN[:ARG1[:ARG2]]  

  Responses:  OK:<TOPIC>  or  ERR:<TOPIC>:<REASON>  

  ─── PING ───  
  PING:PINT                     -> PONG:PONG  

  ─── ON / OFF / TOGGLE ───  
  ON:LAMP                       -> OK:LAMP  
  OFF:LAMP                      -> OK:LAMP  
  TOGGLE:LAMP                   -> OK:LAMP  
  ON:LED                        -> OK:LED  
  OFF:LED                       -> OK:LED  
  TOGGLE:LED                    -> OK:LED  
  ON:BUZZ                       -> OK:BUZZ  
  OFF:BUZZ                      -> OK:BUZZ  

  ─── SET ───  
  SET:LED:MODE:SOLID            -> OK:LED  
  SET:LED:MODE:FADE             -> OK:LED  
  SET:LED:MODE:BLINK            -> OK:LED  
  SET:LED:BRIGHT:<0..255>       -> OK:LED  

  ─── GET ───  
  GET:LAMP:STATE                -> OK:LAMP:STATE:ON/OFF  
  GET:LED:STATE                 -> OK:LED:STATE:ON/OFF  
  GET:LED:MODE                  -> OK:LED:MODE:SOLID/FADE/BLINK  
  GET:LED:BRIGHT                -> OK:LED:BRIGHT:<0..255>  
  GET:UPTIME                    -> OK:UPTIME:<ms>

  ───────────────────────────────────────────────────────────────  
  ▓ FINAL WORDS  
  This is not just a lamp.  
  This is **СоюзПечаль**.  
