#include "oor.hpp"

int* CPPOOR::DetectEllipses(CvMat* imgData) {
    if (imgData == NULL) {
        return NULL;
    }
    Mat src = cvarrToMat(imgData);
    Mat orig_src = src.clone();
    medianBlur(src, src, 3);

    Mat hsv_image;
    cvtColor(src, hsv_image, COLOR_BGR2HSV);

    Mat lower_red_hue_range;
    Mat upper_red_hue_range;
    inRange(hsv_image, Scalar(0,100,100), Scalar(10,255,255), lower_red_hue_range);
    inRange(hsv_image, Scalar(160,100,100), Scalar(179,255,255), upper_red_hue_range);

    Mat red_hue_image;
    addWeighted(lower_red_hue_range, 1.0, upper_red_hue_range, 1.0, 0.0, red_hue_image);

    GaussianBlur(red_hue_image, red_hue_image, Size(9,9), 2, 2);

    vector<vector<Point> > contours;
    vector<Vec4i> hierarchy;
    findContours(red_hue_image, contours, hierarchy, CV_RETR_TREE, CV_CHAIN_APPROX_SIMPLE, Point(0,0));

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
    }

    // paint all red items
    for (int i = 0; i < bounding_rects.size(); i++) {
        Scalar color(0, 255, 0);
        rectangle(orig_src, bounding_rects[i], color, 3, 8, 0);
    }
    // paint largest area blue!
    rectangle(orig_src, bounding_rect, Scalar(255,0,0), 5, 8, 0);

    // calculate center of rectangle
    Point rect_center = (bounding_rect.tl() + bounding_rect.br())*0.5;
    Point orig_src_center = Point(orig_src.size().width/2, orig_src.size().height/2);

    // debug purpose
    // imwrite("result.png", orig_src);

    // create array of size int. Mayby a overkill since we only use 4 indecies.
    int* p = (int *)malloc(sizeof(int));
    p[0] = rect_center.x;
    p[1] = rect_center.y;
    p[2] = orig_src_center.x;
    p[3] = orig_src_center.y;
    
    return p;
}
