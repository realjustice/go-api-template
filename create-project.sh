#!/bin/bash

# Go API Template - 项目创建脚本
# 快速从模板创建新项目

set -e

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
CYAN='\033[0;36m'
NC='\033[0m' # No Color

# 打印函数
print_info() {
    echo -e "${BLUE}ℹ️  $1${NC}"
}

print_success() {
    echo -e "${GREEN}✅ $1${NC}"
}

print_warning() {
    echo -e "${YELLOW}⚠️  $1${NC}"
}

print_error() {
    echo -e "${RED}❌ $1${NC}"
}

print_header() {
    echo ""
    echo -e "${CYAN}========================================${NC}"
    echo -e "${CYAN}   Go API Template - 项目创建工具${NC}"
    echo -e "${CYAN}========================================${NC}"
    echo ""
}

# 检查命令是否存在
check_command() {
    if ! command -v $1 &> /dev/null; then
        print_error "$1 未安装，请先安装"
        exit 1
    fi
}

# 主函数
main() {
    print_header
    
    # 检查必要的命令
    print_info "检查依赖..."
    check_command "go"
    print_success "依赖检查通过"
    
    # 获取模板目录
    TEMPLATE_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
    PARENT_DIR="$(cd "$TEMPLATE_DIR/.." && pwd)"
    
    # 获取项目信息
    echo ""
    print_info "项目名称：用作项目目录名（如：my-api, user-service）"
    read -p "📝 请输入项目名称: " PROJECT_NAME
    
    if [ -z "$PROJECT_NAME" ]; then
        print_error "项目名称不能为空"
        exit 1
    fi
    
    echo ""
    print_info "Go 模块路径：go.mod 中的 module 声明，用于代码 import"
    echo ""
    echo "   常见格式："
    echo "   - github.com/username/$PROJECT_NAME    (开源项目)"
    echo "   - $PROJECT_NAME                         (简单名称，推荐)"
    echo "   - company.com/$PROJECT_NAME            (企业项目)"
    echo ""
    print_info "💡 脚本会自动处理域名格式，无需担心网络下载问题"
    echo ""
    read -p "📦 请输入 Go 模块路径: " MODULE_PATH
    
    if [ -z "$MODULE_PATH" ]; then
        print_error "模块路径不能为空"
        exit 1
    fi
    
    # 去除末尾的斜杠
    MODULE_PATH="${MODULE_PATH%/}"
    
    # 项目完整路径
    PROJECT_PATH="$PARENT_DIR/$PROJECT_NAME"
    
    # 检查项目目录是否已存在
    if [ -d "$PROJECT_PATH" ]; then
        print_error "目录 $PROJECT_PATH 已存在"
        exit 1
    fi
    
    # 确认信息
    echo ""
    print_info "项目信息确认："
    echo "   项目名称: $PROJECT_NAME"
    echo "   模块路径: $MODULE_PATH"
    echo "   创建位置: $PROJECT_PATH"
    echo ""
    read -p "确认创建？(y/N): " CONFIRM
    
    if [[ ! $CONFIRM =~ ^[Yy]$ ]]; then
        print_warning "已取消创建"
        exit 0
    fi
    
    # 创建项目
    echo ""
    print_info "创建项目目录..."
    mkdir -p "$PROJECT_PATH"
    cd "$PROJECT_PATH"
    
    # 创建目录结构
    print_info "创建目录结构..."
    mkdir -p cmd/server
    mkdir -p internal/{controller,service,repository,model,middleware,router,constants}
    mkdir -p pkg/{config,database,redis,cache,logger,errors,web,security,tools}
    mkdir -p config
    mkdir -p logs
    mkdir -p bin
    
    # 复制文件
    print_info "复制模板文件..."
    
    # 复制 pkg 目录
    if [ -d "$TEMPLATE_DIR/pkg" ]; then
        for dir in "$TEMPLATE_DIR/pkg"/*; do
            if [ -d "$dir" ]; then
                dirname=$(basename "$dir")
                cp -r "$dir"/* "./pkg/$dirname/" 2>/dev/null || true
            fi
        done
    fi
    
    # 复制 internal 目录
    if [ -d "$TEMPLATE_DIR/internal" ]; then
        for dir in "$TEMPLATE_DIR/internal"/*; do
            if [ -d "$dir" ]; then
                dirname=$(basename "$dir")
                cp -r "$dir"/* "./internal/$dirname/" 2>/dev/null || true
            fi
        done
    fi
    
    # 复制配置文件
    [ -f "$TEMPLATE_DIR/config/config.yaml" ] && cp "$TEMPLATE_DIR/config/config.yaml" ./config/
    [ -f "$TEMPLATE_DIR/Makefile" ] && cp "$TEMPLATE_DIR/Makefile" ./
    [ -f "$TEMPLATE_DIR/.gitignore" ] && cp "$TEMPLATE_DIR/.gitignore" ./
    
    # 复制 cmd/server 下的入口与 wire 文件（保证与模板一致且参与后续路径替换）
    [ -f "$TEMPLATE_DIR/cmd/server/main.go" ] && cp "$TEMPLATE_DIR/cmd/server/main.go" ./cmd/server/
    [ -f "$TEMPLATE_DIR/cmd/server/wire.go" ] && cp "$TEMPLATE_DIR/cmd/server/wire.go" ./cmd/server/
    
    # 复制 go.mod（保持依赖版本一致）
    print_info "复制依赖配置..."
    if [ -f "$TEMPLATE_DIR/go.mod" ]; then
        cp "$TEMPLATE_DIR/go.mod" ./go.mod
    fi
    # 注意：不复制 go.sum，让 go mod tidy 重新生成
    
    # 修改 go.mod 中的模块名（从 go-api-template 改为临时模块名）
    TEMP_MODULE="golinks-api-template"
    print_info "初始化 Go 模块..."
    
    if [ -f "go.mod" ]; then
        # 替换模块名
        sed -i '' "s|module go-api-template|module $TEMP_MODULE|g" go.mod 2>/dev/null || \
        sed -i "s|module go-api-template|module $TEMP_MODULE|g" go.mod
    else
        # 如果模板没有 go.mod，则创建
        go mod init "$TEMP_MODULE"
    fi
    
    # 替换导入路径为临时模块名（含已复制的 cmd/server/main.go、wire.go）
    print_info "配置导入路径..."
    find . -type f -name "*.go" -exec sed -i '' "s|go-api-template|$TEMP_MODULE|g" {} + 2>/dev/null || \
    find . -type f -name "*.go" -exec sed -i "s|go-api-template|$TEMP_MODULE|g" {} +
    
    # 安装依赖（使用临时模块名，不会触发网络下载）
    print_info "下载依赖..."
    go mod download
    
    print_info "整理依赖..."
    go mod tidy
    
    # 现在替换为真实的模块路径
    if [ "$TEMP_MODULE" != "$MODULE_PATH" ]; then
        print_info "更新模块路径: $MODULE_PATH"
        
        # 替换 go.mod
        sed -i '' "s|module $TEMP_MODULE|module $MODULE_PATH|g" go.mod 2>/dev/null || \
        sed -i "s|module $TEMP_MODULE|module $MODULE_PATH|g" go.mod
        
        # 替换所有 .go 文件中的导入路径
        find . -type f -name "*.go" -exec sed -i '' "s|$TEMP_MODULE|$MODULE_PATH|g" {} + 2>/dev/null || \
        find . -type f -name "*.go" -exec sed -i "s|$TEMP_MODULE|$MODULE_PATH|g" {} +
        
        # 显式替换 cmd/server 下的 main.go、wire.go，确保引用路径被正确替换
        for f in cmd/server/main.go cmd/server/wire.go; do
            if [ -f "$f" ]; then
                sed -i '' "s|$TEMP_MODULE|$MODULE_PATH|g" "$f" 2>/dev/null || sed -i "s|$TEMP_MODULE|$MODULE_PATH|g" "$f"
            fi
        done
        
        print_success "模块路径已更新"
    fi
    
    # 如果是域名格式，添加 replace 指令
    if [[ "$MODULE_PATH" =~ \. ]]; then
        print_info "检测到域名格式，添加本地模块配置..."
        
        cat >> go.mod << EOF

// 本地模块配置：告诉 Go 这是本地模块，不要从网络下载
replace $MODULE_PATH => ./
EOF
        
        print_success "已添加 replace 指令"
    fi
    
    # 复制文档
    print_info "创建项目文档..."
    [ -f "$TEMPLATE_DIR/README.md" ] && cp "$TEMPLATE_DIR/README.md" ./
    [ -f "$TEMPLATE_DIR/STRUCTURE.md" ] && cp "$TEMPLATE_DIR/STRUCTURE.md" ./
    
    # 更新文档中的模块名
    if [ -f "README.md" ]; then
        sed -i '' "s|go-api-template|$MODULE_PATH|g" README.md 2>/dev/null || \
        sed -i "s|go-api-template|$MODULE_PATH|g" README.md
        
        sed -i '' "s|Go API Template|$PROJECT_NAME|g" README.md 2>/dev/null || \
        sed -i "s|Go API Template|$PROJECT_NAME|g" README.md
    fi
    
    if [ -f "STRUCTURE.md" ]; then
        sed -i '' "s|go-api-template|$MODULE_PATH|g" STRUCTURE.md 2>/dev/null || \
        sed -i "s|go-api-template|$MODULE_PATH|g" STRUCTURE.md
    fi
    
    # 完成
    echo ""
    print_success "项目创建成功！"
    echo ""
    echo -e "${CYAN}========================================${NC}"
    echo -e "${GREEN}📁 项目位置: ${NC}$PROJECT_PATH"
    echo -e "${CYAN}========================================${NC}"
    echo ""
    print_info "下一步操作："
    echo ""
    echo "   1. 进入项目目录"
    echo -e "      ${CYAN}cd $PROJECT_PATH${NC}"
    echo ""
    echo "   2. 安装 Wire 工具（如果未安装）"
    echo -e "      ${CYAN}go install github.com/google/wire/cmd/wire@latest${NC}"
    echo ""
    echo "   3. 生成依赖注入代码"
    echo -e "      ${CYAN}make wire${NC}"
    echo ""
    echo "   4. 运行项目"
    echo -e "      ${CYAN}make run${NC}"
    echo ""
    echo "   5. 测试 API"
    echo -e "      ${CYAN}curl http://localhost:8080/health${NC}"
    echo ""
    echo -e "${CYAN}========================================${NC}"
    print_success "开始构建你的 API 项目！🚀"
    echo ""
}

# 运行主函数
main
