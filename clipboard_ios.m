// Copyright 2021 The golang.design Initiative Authors.
// All rights reserved. Use of this source code is governed
// by a MIT license that can be found in the LICENSE file.
//
// Written by Changkun Ou <changkun.de>

//go:build ios

#import <UIKit/UIKit.h>
#import <MobileCoreServices/MobileCoreServices.h>

void clipboard_write_string(char *s) {
    NSString *value = [NSString stringWithUTF8String:s];
    [[UIPasteboard generalPasteboard] setString:value];
}

char *clipboard_read_string() {
    NSString *str = [[UIPasteboard generalPasteboard] string];
    return (char *)[str UTF8String];
}

unsigned int clipboard_read_image(void **out) {
    UIPasteboard *pasteboard = [UIPasteboard generalPasteboard];
    UIImage *image = pasteboard.image;
    if (image == nil) {
        return 0;
    }
    NSData *data = UIImagePNGRepresentation(image);
    if (data == nil) {
        return 0;
    }
    NSUInteger siz = [data length];
    *out = malloc(siz);
    [data getBytes: *out length: siz];
    return siz;
}

unsigned int clipboard_read_image_jpeg(void **out) {
    UIPasteboard *pasteboard = [UIPasteboard generalPasteboard];
    UIImage *image = pasteboard.image;
    if (image == nil) {
        return 0;
    }
    NSData *jpegData = UIImageJPEGRepresentation(image, 1.0);
    if (jpegData == nil) {
        return 0;
    }
    NSUInteger siz = [jpegData length];
    *out = malloc(siz);
    [jpegData getBytes: *out length: siz];
    return siz;
}

int clipboard_write_image(const void *bytes, unsigned long n) {
    UIPasteboard *pasteboard = [UIPasteboard generalPasteboard];
    NSData *data = [NSData dataWithBytes: bytes length: n];
    UIImage *image = [UIImage imageWithData:data];
    if (image == nil) {
        return -1;
    }
    pasteboard.image = image;
    return 0;
}

int clipboard_write_image_jpeg(const void *bytes, unsigned long n) {
    UIPasteboard *pasteboard = [UIPasteboard generalPasteboard];
    NSData *jpegData = [NSData dataWithBytes: bytes length: n];
    UIImage *image = [UIImage imageWithData:jpegData];
    if (image == nil) {
        return -1;
    }
    pasteboard.image = image;
    return 0;
}
