package main

// layer 구분을 통한 이미지 압축(코드와 실행파일만 존재)

import "carte_cli/cmd"

func main() {
    cmd.Execute()
}


// =============================여기에 Carte 코드 실행==============================

// terminal실행
// 1) 실행 파일을 시스템 PATH에 추가


// [[[[[ 코드 실행 ]]]]]
// mkdir -p ~/bin
// cp Carte ~/bin/
// echo 'export PATH=$PATH:~/bin' >> ~/.bashrc
// source ~/.bashrc
// [권한 부여]
// sudo mkdir /Carte
// sudo mkdir /Carte/images

// whoami
// sudo chown -R yj /Carte/images


// ==============================================================================
// 2) '~bin'디렉토리가 시스템 PATH에 포함되어있는가 확인하고 포함되지 않은경우 '.bashrc', '.zshrc', '.profile'에 추가
// echo 'export PATH=$PATH:~/bin' >> ~/.bashrc
// source ~/.bashrc


// 2_2) 'zsh'사용하는 경우 '.zshrc'에 파일 추가
// echo 'export PATH=$PATH:~/bin' >> ~/.zshrc
// source ~/.zshrc
