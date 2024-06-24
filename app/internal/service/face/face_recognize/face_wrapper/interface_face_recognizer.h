#pragma once

#ifdef __cplusplus
extern "C" {
#endif

        // images cluster
        /// Image info
        #define MAX_FILE_NAME_LEN 256
        typedef struct
        {
            char         filename[MAX_FILE_NAME_LEN];          ///< filename of image
            float        similarity;                           ///< simility between curent image and query image.
        }ImageInfo;


        enum ImageDataType
        {
            Unknown = 0,                           // unknown, not a valid image data type.
            JPGFile = 1,                           // the character stream of jpg file.
            PNGFile = 2,                           // the character stream of png file.
            NV12Data = 3,                          // Planar YUV420 with interleaved UV, each channel is 8-bit.
            NV21Data = 4,                          // Planar YUV420 with interleaved VU, each channel is 8-bit.
            YUV420Data = 5,                        // Planar YUV420, each channel is 8-bit.
            RGBData = 6,                           // interleaved RGB (in order in memory), each channel is 8-bit.
            BGRData = 7                            // interleaved BGR (in order in memory), each channel is 8-bit.
        };

        typedef struct
        {
            unsigned char* data;                   // the character stream of image file or rows array of image.
            int data_len;                          // the length of data.
            int width;                             // image width, if type is PNGFile or JPGFile, the width is 0.
            int height;                            // image height, if type is PNGFile or JPGFile, the height is 0.
            enum ImageDataType data_type;               // the type of data.
        } ImageData;


        /**
        Initialize model
        @param conf_thresh  similarity.
        @param topk
        @param model_path
        */
        int hiarClusterInit(const float conf_thresh, const int top_k, const char* model_path, const char* logger_path, const char* db_path);
        /**
        adding images
        @param input_list      [input] list of image filename which will be registed.fg.["D:/hiar/gallery/1.jpg", "D:/hiar/gallery/2.jpg"]
        @param len_input       [input] the lenth of input_list array.
        @param ok_list         [output] the buffer of adding successful images.
        @param ok_ok           [input] the length of ok_list buffer,
        @return  int           the number of adding successful images.
        */
        int hiarAddingImages(const ImageInfo* input_list, const int len_input, ImageInfo* ok_list, const int len_ok);
        /**
        adding images
        @param image        the inputing image include more than one faces.
        @param filename     the file name of the image.
        */
        int hiarAddingImage(const ImageData* image,const char* filename);

        /**
        query image
        @param image        the inputing image include more than one faces.
        @param vquery_images  query results.
        */
        int hiarQuery(const ImageData *image, ImageInfo *vquery_images, int v_len );

         /**
        delete no use image
        @param vdel_images        the inputing images, if vdel_images is empty or v_len<=0, all images that added are deleted.
        */
        int hiarDelImages(const ImageInfo* vdel_images, const int v_len);


#ifdef __cplusplus
}
#endif
