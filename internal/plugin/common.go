package plugin

import "sync"

var once sync.Once // 确保初始化逻辑只执行一次
