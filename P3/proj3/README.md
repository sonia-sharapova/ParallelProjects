
#### NOTE:
I originally wrote my implementation on my local device and faced unexpected hurdles when trying to install OpenCV on the Linux Cluster. All libraries had to be installed manually and some errors were unresolved. 

### Requirements
Required Python packages
pip install pandas matplotlib seaborn


#### For the go program:

To Run:
Initialize go module:
    go mod init proj3
Install required dependencies:
	go get -u gocv.io/x/gocv 
	go get -u github.com/suyashkumar/dicom 
	go mod tidy

For cluster:
    sbatch ./benchmark/run_benchmark.sh

Independent Runs:
- Sequential:
    $ go run ./cmdprocess-dicom/main.go -input ./smallerData -output ./output/sequential/ -mode sequential $

- Pipeline:
    $ go run ./cmd/process-dicom/main.go -input ./smallerData -output ./output/pipeline/ -mode pipeline -workers 8 -buffer 10 

- Work Stealing:
    $ go run ./cmd/process-dicom/main.go -input ./smallerData -output ./output/workstealing/ -mode workstealing -workers 8

Run All (Benchmarking)
    $ go run ./cmd/benchmark/main.go -input ./smallerData -maxworkers 8


#### To get OpenCV on Linux:

I included the tar zip file in the repository.
opencv_contrib has to be taken from https://github.com/opencv/opencv_contrib/tree/4.x
- move into the opencv directory

Note: make sure the paths for modules are right!

    $ cd ~/opencv/opencv-4.10.0
    $ mkdir -p build && cd build
    $ cmake -DCMAKE_MODULE_PATH=~/opencv/opencv_contrib/modules \
        -DCMAKE_BUILD_TYPE=Release \
        -DCMAKE_INSTALL_PREFIX=/home/sharapova/opencv_local \
        -DBUILD_opencv_python3=ON \
        -DBUILD_opencv_python2=OFF \
        -DBUILD_TESTS=OFF \
        -DBUILD_EXAMPLES=OFF \
        -DWITH_TBB=ON \
        -DWITH_CUDA=OFF \
        -DWITH_OPENGL=ON \
        -DWITH_FFMPEG=ON \
        -DBUILD_SHARED_LIBS=ON \
        -DPROTOBUF_UPDATE_FILES=ON \
        -DBUILD_PROTOBUF=OFF \
        -DOPENCV_EXTRA_MODULES_PATH=~/opencv/opencv_contrib/modules \
        -DBUILD_LIST=core,imgproc,highgui,imgcodecs,aruco,dnn,video,videoio,calib3d,features2d,flann,ml,objdetect,photo,stitching,tracking,ts \
        ..
    $ make -j$(nproc)
    $ make install


add the following to your .bashrc file:
    export CGO_CXXFLAGS="-I/home/sharapova/opencv_local/include"
    export CGO_CPPFLAGS="-I/home/sharapova/opencv_local/include/opencv4"
    export CGO_LDFLAGS="-L/home/sharapova/opencv_local/lib -lopencv_core -lopencv_imgproc -lopencv_highgui -lopencv_imgcodecs -lopencv_aruco"
    export LD_LIBRARY_PATH="/home/sharapova/opencv_local/lib:$LD_LIBRARY_PATH"
    export PKG_CONFIG_PATH="/home/sharapova/opencv_local/lib/pkgconfig"