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

- [avv help](#avv-help)
- [avv make](#avv-make)
- [avv run](#avv-run)
- [avv version](#avv-version)

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

### avv run 
作ったタスクを実行します

```sh
Usage:
  avv run [flags]

Flags:
      --all
  -n, --count int              number of task (default 1)
  -y, --extractParallel int    WaveViewの並列数です (default 1)
  -h, --help                   help for run
  -m, --maxRetry int           各ステージの処理が失敗したときに再実行する回数です
 (default 3)
  -x, --simulateParallel int   HSPICEの並列数です (default 1)
      --slack                  すべてのタスクが終わったときにSlackへ投稿します (d
efault true)

Global Flags:
      --config string   config file (default is $HOME/.config/avv/.avv.json)
```

#### avv run all, n,count オプション
一度に複数のタスクを実行できます。個数を個別に指定する場合は `n,count` オプション。保存されているタスクをすべて実行するには `all` オプションを指定します

- `--all`       --> すべてのタスクを実行します
- `-n, --count n` --> n個のタスクを実行します。作られた順番が保存されないことに注意してください

```sh
# ex) 8個実行
$ avv run -n 8


# ex) すべて
$ avv run --all
```

#### avv run Parallelオプション
`avv` では Hspice と WaveView を それぞれ並列で動かしています

```
                               / => [Hspice] => \
[First task pool]-- pop one => - => [Hspice] => - => [simulation result] => [Second task pool]
                               \ => [Hspice] => /


                                / => [WaveView] => \
[Second task pool]-- pop one => - => [WaveView] => - => [extract result] => [csv]
                                \ => [WaveView] => /

```

Parallel オプションを使うと、HspiceやWaveViewの並列数を決めることが可能です

-  `-x, --simulateParallel int` -->  HSPICEの並列数です (default 1)
-  `-y, --extractParallel int`  -->  WaveViewの並列数です (default 1)

`avv` の性能を決める一番重要なオプションであることに注意してください。例えば50,000回のモンテカルロ・シミュレーションは Hspice の段階で負荷になることはあまりありませんが、WaveViewではかなり多くのメモリが必要になります。このためWaveViewの並列数を大きくしすぎるとハングアップすることがあります。動作させるOCのスペックと相談する必要があります

```sh
# ex) HSpiceが10並列 WaveViewが3並列
$ avv run -x 10 -y 3
```

#### avv run maxRetry オプション
`avv` ではHspice, WaveViewが何らかの原因で失敗した場合、maxRetryに指定した回数だけ再実行を試みます。

```sh
# ex) 再実行回数の上限が5
$ avv run --maxRetry 5
```

#### avv run slack オプション
`avv` はタスクがすべて終了したとき、slackへ投稿することができます。詳細は[SlackConfig](#SlackConfig) を参照してください。デフォルトでTrueになります

- `--slack`  -->  すべてのタスクが終わったときにSlackへ投稿します (dfault true)

```sh
# ex) slackへ投稿しない
$ avv run --slack=false
```

### avv version
バージョン情報を出力して終了します

```sh
# ex) バージョン情報を出力する
$ avv version
```

## Config
`avv` のコンフィグについて記述します。これらを正確に設定しないと、目的のシミュレーションが実行できないことに注意してください

| 項目 | 値 |  
|:----:|:--:|  
|場所|`$HOME/.config/avv/.avv.json`|  
|形式|JSON|  

[example](./example/avv.json)

### Templates
#### SPIScript
シミュレーションしたい回路のSPIスクリプトのテンプレートです

1. 素のnetlistを用意します。多くの場合 `~/simulation/username/name/name/HSPICE/nominal/netlist` がそれです。加工するので、コピーをどこかに作成します
2. コピーしたファイルの先頭に以下の行を追加します

```
.option MCBRIEF=2
.param vtn=AGAUSS(%.4f,%.4f,%.4f) vtp=AGAUSS(%.4f,%.4f,%.4f)
.option PARHIER = LOCAL
.include '%s'
.option ARTIST=2 PSF=2
.temp 25
.include '%s'
```

3. 更に最後尾に以下の行を追加します

```
.tran 10p 20n start=0 uic sweep monte=%d firstrun=1
.option opfile=1 split_dp=2

.end
```

4. 作成したファイルへのパスを `SPIScript` の値として設定します

### Default
シミュレーションに必要な各パラメータのデフォルト値を決めます。コマンドラインオプションから指定できるものだけでなくここでしか設定できないものがあります。

|項目|説明|CLIオプション|例|  
|:--:|:--:|:--:|:--:|  
|`NetListDir`|`avv` が生成するSPIスクリプトを設置する場所|なし|`$HOME/simulation/username/name/name/monte_carlo/netlist`|  
|`BaseDir`|`avv` がワークスペースとして使うディレクトリのルートです|avv make --basedir|`$HOME/Workspace/base`|  
|SEED.Start|Seed値の開始値です|avv make --start|1|  
|SEED.End|Seed値の終了値です|avv make --end|2000|  
|Parameters.PlotPoint.Start|プロットの始点です|avv make -a, --PlotStart|2.5E-9|  
|Parameters.PlotPoint.Step|プロットの刻み幅です|avv make -b, --PlotStep|7.5E-9|  
|Parameters.PlotPoint.Stop|プロットの終点です|avv make -c, --PlotStop|1.75E-8|  
|Parameters.PlotPoint.Signals|取り出す信号線のリストです|なし|["N1","N2","BLB"]|  
|Parameters.Sweeps|一回あたりのモンテカルロ・シミュレーションの回数です|avv make -t, --times|5000|  
|Parameters.Vtn.Threshold|Vtnのしきい値電圧です|avv make -N, --VtnVoltage|0.6|  
|Parameters.Vtn.Sigma|Vtnのシグマです|avv make --VtnSigma|0.046|  
|Parameters.Vtn.Deviation|Vtnの偏差です|avv make --VtnDeviation|1.0|  
|Parameters.Vtp.Threshold|Vtpのしきい値電圧です|avv make -N, --VtpVoltage|-0.6|  
|Parameters.Vtp.Sigma|Vtpのシグマです|avv make --VtpSigma|0.046|  
|Parameters.Vtp.Deviation|Vtpの偏差です|avv make --VtpDeviation|1.0|  
|Parameters.AddFile.VddVoltage|電源電圧です|なし|0.8|  
|Parameters.AddFile.GndVoltage|グランドの電圧です|なし|0|  
|Parameters.AddFile.ICCommand|N1やN2などの初期値を決める .IC コマンドを記述します|なし|".IC V(N1)=0.8V V(N2)=0V"|  
|Parameters.AddFile.Options|その他のオプションがあれば記述します|なし|[]|  
|Parameters.ModelFile|モデルファイルへの __絶対パス__ です|なし|`$HOME/Workspace/avv/modelfile.exe`|  
### SlackConfig
Slackへ投稿する際に必要な設定です

|項目|説明|CLIオプション|例|  
|:--:|:--:|:--:|:--:|  
|SlackConfig.User|@付きでメンションする相手です。自分にしておくといいと思います|なし|"xztaityozx"|  
|SlackConfig.WebHookURL|SlackのWebHookURLです。ここでは解説しません。ggってください。この値は外部に公開しないでください|なし|URL|  
|SlackConfig.Channel|投稿するチャンネルです|なし|"#進捗"|  
|SlackConfig.MachineName|動作しているマシンに任意の名前をつけられます。複数のマシンで動いている `avv` を識別するために使います|なし|"MachineA"|  

### HSPICE
Hspiceの設定です

|項目|説明|CLIオプション|例|  
|:--:|:--:|:--:|:--:|  
|HSPICE.Path|Hspiceの実行ファイルへの __絶対パス__ です|なし|"/path/to/hspice"|  
|HSPICE.Options|Hspiceへ渡すオプションです|なし|""|  

### WaveView
WaveViewの設定です

|項目|説明|CLIオプション|例|  
|:--:|:--:|:--:|:--:|  
|WaveView.Path|WaveViewの実行ファイルへの __絶対パス__ です|なし|"/path/to/wv"|  
