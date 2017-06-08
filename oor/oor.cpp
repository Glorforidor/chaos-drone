#include "oor.hpp"

bool CompareTolerance(int a, int b, int t) {
    return abs(a - b) <= t;
}

bool CompareTolerance(RotatedRect a, RotatedRect b, int t) {
    return CompareTolerance(a.center.x, b.center.x, 5) && CompareTolerance(a.center.y, b.center.y, 5) &&
        CompareTolerance(a.size.width, b.size.width, t) && CompareTolerance(a.size.height, b.size.height, t);
}

int* CPPOOR::DetectEllipses(CvMat* imgData) {
    if (imgData == NULL) {
        return NULL;
    }
    Mat src = cvarrToMat(imgData);
    Mat orig_src = src.clone();
    // medianBlur(src, src, 1);

    Mat hsv_image;
    cvtColor(src, hsv_image, COLOR_BGR2HSV);
    // imwrite("hsv.png", src);

    Mat lower_red_hue_range;
    Mat upper_red_hue_range;
    inRange(hsv_image, Scalar(0,50,80), Scalar(10,255,255), lower_red_hue_range);
    inRange(hsv_image, Scalar(140,40,50), Scalar(180,255,255), upper_red_hue_range);

    Mat red_hue_image;
    addWeighted(lower_red_hue_range, 1.0, upper_red_hue_range, 1.0, 0.0, red_hue_image);

    GaussianBlur(red_hue_image, red_hue_image, Size(7,7), 0, 0);

    // imwrite("gauss.png", red_hue_image);

    // vector<cv::Vec3f> circles;
    // HoughCircles(red_hue_image, circles, CV_HOUGH_GRADIENT, 1, red_hue_image.rows/8, 100, 20, 0, 0);
    // for(size_t current_circle = 0; current_circle < circles.size(); ++current_circle) {
    // Point center(round(circles[current_circle][0]), round(circles[current_circle][1]));
    // int radius = round(circles[current_circle][2]);
    // circle(orig_src, center, radius, cv::Scalar(0, 255, 0), 2);
    // }

    // Detects canny edges for the image, which we will use to compare to our
    // approach of using findContours.
    Mat edge_image = red_hue_image.clone();
    Canny(edge_image, edge_image, 80.0, 100.0);
    imwrite("edgeimg.png", edge_image);

    vector<vector<Point> > contours;
    // vector<Vec4i> hierarchy;
    findContours(red_hue_image, contours, CV_RETR_TREE, CV_CHAIN_APPROX_SIMPLE);

    vector<RotatedRect> minEllipse( contours.size() );

    double imgWidth = orig_src.size().width;
    double imgHeight = orig_src.size().height;

    int largest_area = 0;
    RotatedRect largest_ellipse_fit;
    int largest_area_index;
    Rect bounding_rect;
    //vector<Rect> bounding_rects;
    for (int i = 0; i < contours.size(); i++) {
        double a = contourArea(contours[i], false);
        if (a >= largest_area) {
            Rect tmp_bounding_rect = boundingRect(contours[i]);
            largest_area = a;
            if (contours[i].size() > 50 && tmp_bounding_rect.x > 0 && tmp_bounding_rect.y > 0 && tmp_bounding_rect.x + tmp_bounding_rect.width < imgWidth &&
                tmp_bounding_rect.y + tmp_bounding_rect.height < imgHeight) {
                bounding_rect = tmp_bounding_rect;
                //bounding_rects.push_back(bounding_rect);
                largest_ellipse_fit = fitEllipse(Mat(contours[i]));
                largest_area_index = i;
            }
        }
        if( contours[i].size() > 60 ) {
            minEllipse[i] = fitEllipse( Mat(contours[i]) );
            //cout << minEllipse[i].center << endl;
            ellipse(orig_src, minEllipse[i], Scalar(0, 255, 0), 2, 8);
        }
    }
    bool isValid;
    for (int i = 0; i < contours.size(); i++) {
        if (i != largest_area_index && contours[i].size() > 50) {
            RotatedRect tmp_ellipse = fitEllipse(Mat(contours[i]));
            if (CompareTolerance(largest_ellipse_fit, tmp_ellipse, 50)) {
                isValid = true;
                ellipse(orig_src, largest_ellipse_fit, Scalar(255, 0, 0), 2, 8);
                break;
            }
        }
    }
    imwrite("aftercontours.png", red_hue_image);
    imwrite("fit.png", orig_src);

    // paint all red items
    // for (int i = 0; i < bounding_rects.size(); i++) {
    // Scalar color(0, 255, 0);
    // rectangle(orig_src, bounding_rects[i], color, 1, 8, 0);
    // }
    // paint largest area blue!
    rectangle(orig_src, bounding_rect, Scalar(255,0,0), 1, 8, 0);

    // calculate center of rectangle
    Point orig_src_center = Point(imgWidth/2, imgHeight/2);

    // debug purpose
    //imwrite("result.png", orig_src);
    //imwrite("red.png", red_hue_image);

    if (!isValid || bounding_rect.width < 70 || bounding_rect.height < 70) {
        bounding_rect.x = 0;
        bounding_rect.y = 0;
        bounding_rect.width = 0;
        bounding_rect.height = 0;
    }

    // create array of size int. Mayby a overkill since we only use 4 indecies.
    int* p = (int *)malloc(sizeof(int));
    p[0] = bounding_rect.x;
    p[1] = bounding_rect.y;
    p[2] = bounding_rect.width;
    p[3] = bounding_rect.height;
    p[4] = orig_src_center.x;
    p[5] = orig_src_center.y;

    return p;
}
