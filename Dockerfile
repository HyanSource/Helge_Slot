# 鏡像檔
FROM golang:alpine
# 工作目錄
WORKDIR /server
# 加入檔案
ADD . /server
# 進入資料夾和建構執行檔案
RUN cd /server && go build Server.go
# 阜號
EXPOSE 8124
# 執行檔案
ENTRYPOINT ./Server