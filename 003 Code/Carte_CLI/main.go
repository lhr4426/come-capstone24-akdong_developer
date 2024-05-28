package main

import "carte_cli/cmd"

func main() {
    cmd.Execute()
}


// terminal실행
// 1) 실행 파일을 시스템 PATH에 추가
// mkdir -p ~/bin
// cp Carte ~/bin/

// 2) '~bin'디렉토리가 시스템 PATH에 포함되어있는가 확인하고 포함되지 않은경우 '.bashrc', '.zshrc', '.profile'에 추가
// echo 'export PATH=$PATH:~/bin' >> ~/.bashrc
// source ~/.bashrc


// 2_2) 'zsh'사용하는 경우 '.zshrc'에 파일 추가
// echo 'export PATH=$PATH:~/bin' >> ~/.zshrc
// source ~/.zshrc
