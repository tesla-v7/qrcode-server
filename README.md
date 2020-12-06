# Generation of one-time qr codes.

### Instaling

1) Clone project
```
git clone https://github.com/tesla-v7/qrcode-server.git
```
2) Change folder
```
cd  qrcode-server
```
3) Configure your config.toml
```
[qr]

#maximum identifier value
idMax=999

#prefix for idMax(required for deployment on multiple servers)
idPrefix=1

#logo file name
logoPath="logo.png"

#color in the center of qr code
colorCenter="05ffb8"

#qr code edge color
colorEdge="007aff"

#radius of one bit of qr code
pixRadius=5

#gr code lifetime
lifetime=120

#buffer size of ready qr codes
sizeBuffer=8

#the number of threads to generate qr codes
numberOfThreads=1

#coded text pattern in qr code
template="{\"id\": %d}"
```

#Build
```
./build.sh
```

# Run
```
./qr-code s --listen 0.0.0.0:3344
```

# View result
```
curl http://0.0.0.0:3344/qrCode
```
Result
```
{"id":9281,"qrBase64":"data:image/png;base64,iVBOR...kJggg=="}
```
qr-code result image:
 
![qrResult.png](qrResult.png) 

## Authors

* **Vitalii Vidanov** - [Vitalii](https://github.com/tesla-v7)

## License

This project is licensed under the MIT License - see the [LICENSE.md](LICENSE.md) file for details