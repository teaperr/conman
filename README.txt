  ___   ___   _ __   _ __ ___    __ _  _ __  
 / __| / _ \ | '_ \ | '_ \ _ \  / _' || '_ \ 
| (__ | (_) || | | || | | | | || (_| || | | |
 \___| \___/ |_| |_||_| |_| |_| \__,_||_| |_|
                                             
         a (con)figuration (man)ager

what the hell is this?
  conman is a simple program designed to:
    - help you out with annoying sprawling configs
    - make your configurations more portable
    - make configuring multiple systems more "plug and play" than wasting time on moving files around

why did i make conman?
  i got a thinkpad recently and wanted to easily sync my configurations using syncthing without needing to add each file manually every time.

how do i install?
  dependencies:
    go
    gcc-go
  
  arch using yay:
      yay -S conman```
  other distros:
    git clone https://github.com/teaperr/conman
    cd conman
    make
    sudo make install
    cd ../
    rm -r conman/

how do i support development?
  suggestions for features, bug reports, pull requests etc are all greatly appreciated!
  
  you can also support me via my paypal at https://paypal.me/teaperr or my ko-fi at https://ko-fi.com/teaper_
  (paypal has less fees)
