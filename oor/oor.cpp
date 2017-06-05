#include "oor.hpp"

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

    GaussianBlur(red_hue_image, red_hue_image, Size(7,7), 3, 3);

    // imwrite("gauss.png", red_hue_image);

    // vector<cv::Vec3f> circles;
    // HoughCircles(red_hue_image, circles, CV_HOUGH_GRADIENT, 1, red_hue_image.rows/8, 100, 20, 0, 0);
    // for(size_t current_circle = 0; current_circle < circles.size(); ++current_circle) {
        // Point center(round(circles[current_circle][0]), round(circles[current_circle][1]));
        // int radius = round(circles[current_circle][2]);
        // circle(orig_src, center, radius, cv::Scalar(0, 255, 0), 2);
    // }



    vector<vector<Point> > contours;
    // vector<Vec4i> hierarchy;
    findContours(red_hue_image, contours, CV_RETR_TREE, CV_CHAIN_APPROX_SIMPLE);

    vector<RotatedRect> minEllipse( contours.size() );


    int largest_area = 0;
    Rect bounding_rect;
    vector<Rect> bounding_rects;
    for (int i = 0; i < contours.size(); i++) {
        double a = contourArea(contours[i], false);
        if (a >= largest_area) {
            largest_area = a;
            bounding_rect = boundingRect(contours[i]);
            bounding_rects.push_back(bounding_rect);
        }
        if( contours[i].size() > 60 ) {
            minEllipse[i] = fitEllipse( Mat(contours[i]) );
            cout << minEllipse[i].center << endl;
            ellipse(orig_src, minEllipse[i], Scalar(0, 255, 0), 2, 8);
        }
    }

    // paint all red items
    // for (int i = 0; i < bounding_rects.size(); i++) {
        // Scalar color(0, 255, 0);
        // rectangle(orig_src, bounding_rects[i], color, 1, 8, 0);
    // }
    // paint largest area blue!
    rectangle(orig_src, bounding_rect, Scalar(255,0,0), 1, 8, 0);

    // calculate center of rectangle
    Point orig_src_center = Point(orig_src.size().width/2, orig_src.size().height/2);

    // debug purpose
    imwrite("result.png", orig_src);
    imwrite("red.png", red_hue_image);

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
