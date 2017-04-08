#include "oor.hpp"
#include "oor.h"

OOR OORInit()
{
    CPPOOR* ret = new CPPOOR();
    return (void*)ret;
}

void OORFree(OOR obj)
{
    CPPOOR* foo = (CPPOOR*)obj;
    delete foo;
}

int* DetectEllipses(OOR obj, void* imgData)
{
    CPPOOR* foo = (CPPOOR*)obj;
    CvMat* img = (CvMat*)imgData;
    return foo->DetectEllipses(img);
}
