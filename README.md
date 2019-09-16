# avv
avv は HspiceとWaveViewを使ってモンテカルロ・シミュレーションをするツールです

## Requirements
- CentOS 7.5 以上
  - Hspiceの動作環境がRedHat系なため
- git 2.20.1 以上
- Hspice 2017.03-SP2-4 以上
- Custom WaveView 2017.12-SP1-1 以上

### Optional
- GoLang v1.12.5 以上

## Install
### GoLang をインストールした場合
```sh
# github から最新を取得してインストール
$ go get -u github.com/xztaityozx/avv
```

### バイナリをダウンロードする方法
GitHubの[Releaseページ](https://github.com/xztaityozx/avv/releases)からバイナリをダウンロードする

```sh
# v3.2.71をダウンロードしたとする
$ tar xvfz ./avv-v3.2.71.tar.gz
... tar outputs ...

# PATHが通ってるところにバイナリを置く
$ cp ./build/avv ~/.local/bin/
```

## Usage
`avv` は複数のサブコマンドで機能を切り替えて使います

```sh
avv [command]

Available Commands:
  help        Help about any command
  make        タスクを作ります
  run         シミュレーションを実行します
  version     Print avv version

Flags:
      --config string   config file (default is $HOME/.config/avv/.avv.json)
  -h, --help            help for avv

Use "avv [command] --help" for more information about a command.
```

- [avv help](#avv help)
- [avv make](#avv make)
- [avv run](#avv run)
- [avv version](#avv version)

### avv help
コマンドのヘルプを出力します。詳細なリファレンスを見ることなく、おおよその使い方がわかります

```sh
Usage: avv help [command] [flags]

# ex) サブコマンドrunのヘルプを見る
$ avv help run
```

### avv make
タスクを作ります。作ったタスクはJSONファイルとして保存され、タスクが正常に終了したときに削除されます。タスクファイルの保存場所は[Config](#Default)の `BaseDir` に設定されたディレクトリの下、 `task` ディレクトリに保存されます(例えばBaseDirが `/home/user/base` なら、 `/home/user/base/task` です)

`float` と書いてあるところは浮動小数点の値を与えることができます。10のn乗を表す、E表記も解釈することができます。

```
1000 => 1E3
0.1 => 1E-1
```

```sh
Usage:
  avv make [flags]

Flags:
  -a, --PlotStart float      プロットの始点 (default 2.5e-09)
  -b, --PlotStep float       プロットの刻み幅 (default 7.5e-09)
  -c, --PlotStop float       プロットの終点 (default 1.75e-08)
      --VtnDeviation float   Vtnの偏差 (default 1)
      --VtnSigma float       Vtnのシグマ
  -N, --VtnVoltage float     Vtnのしきい値電圧
      --VtpDeviation float   Vtpの偏差 (default 1)
      --VtpSigma float       Vtpのシグマ
  -P, --VtpVoltage float     Vtpのしきい値電圧
      --basedir string       シミュレーションの結果を書き出す親ディレクトリ
      --end int              SEEDの終点 (default 10)
  -h, --help                 help for make
  -S, --sigma float          VtpとVtnのシグマ (default 0.046)
      --start int            SEEDの始点 (default 1)
  -t, --times int            モンテカルロシミュレーション1回当たりの回数 (default 100)

Global Flags:
      --config string   config file (default is $HOME/.config/avv/.avv.json)
```

`(default value)` と記載のあるオプションには指定しなかった場合に自動的に設定されるデフォルト値があります。この記載のないオプションを設定しない場合は[Config](#Default)の値が使われます。値の優先度は以下です

```
コマンドラインオプション > Config.Default (> Default Value)
```

#### avv make transistorオプション
- `-N, --VtnVoltage` --> Vtnのしきい値電圧です
- `--VtnSigma` --> Vtnのシグマです
- `--VtnDeviation` --> Vtnの偏差です
- `-N, --VtpVoltage` --> Vtpのしきい値電圧です
- `--VtpSigma` --> Vtpのシグマです
- `--VtpDeviation` --> Vtpの偏差です
- `-S,--sigma` --> VtnとVtpのシグマを同時に指定できます
  - 個別の設定がある(`--VtnSigma,--VtpSigma` が指定されている)場合はそちらが優先されます

```sh
# ex) vtnがしきい値1, シグマが2, 偏差が3のとき
$ avv make -N 1 --VtnSigma 2 --VtnDeviation 3
# ex) vtnがしきい値0.5, シグマが0.14, 偏差が1でvtpはしきい値が-0.4, シグマが0.14 偏差が1のとき
$ avv make -N 0.5 -P -0.4 --VtnSigma 0.14 --VtpSigma 0.14 --VtnDeviation 1 --VtpDeviation 1
# もしくは
$ avv make -N 0.5 -P -0.4 -S 0.14
```

#### avv make PlotPoint オプション
WaveViewで波形データから、特定の信号を取り出すときに、時間を指定することができます。設定には[開始点],[刻み幅],[終点]指定します

- `-a, --PlotStart float` -->  プロットの始点 (default 2.5e-09)
- `-b, --PlotStep float`  -->  プロットの刻み幅 (default 7.5e-09)
- `-c, --PlotStop float`  -->  プロットの終点 (default 1.75e-08)

```sh
# ex) 1sから2sごとに10sまでの点(1s,3s,5s,...,9s)を指定する
$ avv make -a 1 -b 2 -c 10
```

#### avv make seed オプション
`avv` ではTransistorなどの設定を固定したまま、Seed値を変化させながらタスクを作ることが可能です

- `--start int`   -->     SEEDの始点 (default 1)
- `--end int`     -->     SEEDの終点 (default 10)

```sh
# ex) Seedが1から2000の2000個のタスクを作る
$ avv make --end 2000
```

#### avv make times オプション
モンテカルロ・シミュレーション1回あたりの回数を指定できます

- `-t, --times int` -->  モンテカルロシミュレーション1回当たりの回数 (default

## Config

### Templates
### Default
### SlackConfig
### AutoRemove
### HSPICE
### WaveView
