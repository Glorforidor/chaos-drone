#include <opencv2/core/types_c.h>

#ifdef __cplusplus
extern "C" {
#endif
  typedef void* OOR;
  OOR OORInit(void);
  void OORFree(OOR);
  int* DetectEllipses(OOR, void*);
#ifdef __cplusplus
}
#endif
