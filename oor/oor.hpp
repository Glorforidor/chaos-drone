#include <opencv2/imgproc/imgproc.hpp>
#include <opencv2/highgui/highgui.hpp>
#include <opencv2/features2d/features2d.hpp>
#include <iostream>
#include <stdlib.h>

using namespace cv;
using namespace std;

class CPPOOR
{
public:
    CPPOOR(void){};
    ~CPPOOR(){};
    int* DetectEllipses(CvMat* imgData);
};
