/***************************************************************************
 * 
 * Copyright (c) 2014 Baidu.com, Inc. All Rights Reserved
 * 
 **************************************************************************/
 
 
 
/**
 * @file test.cpp
 * @author quanzongfeng(com@baidu.com)
 * @date 2014/03/13 14:54:03
 * @brief 
 *  
 **/

#include "lbfgs.h"
#include <vector>
#include <string>
#include <iostream>

using namespace std;

double gradient(double x) {
    return 4.0 * x* x*x;
}

double f(double x) {
    return x*x*x*x;
}

int main() {
    double x = 10.0;
    double gx = 0.0;
    double fx = 0.0;
    CRFPP::LBFGS lb = CRFPP::LBFGS();
    
    int r = 0;
    while(1) {
        gx = gradient(x);
        fx = f(x);
        cout << x << "\t"<<gx<<"\t"<<fx<<endl;
        r = lb.optimize(1, &x, fx, &gx, false, 1);
        if (r < 0) {
            cout << "optimzie error "<<r<<endl;
            break;
        }
        if (r == 0) {
            cout << "endl with x:=" << x<<endl;
            break;
        }
    }
    return 1;
}
    

        

        

        























/* vim: set expandtab ts=4 sw=4 sts=4 tw=100: */
