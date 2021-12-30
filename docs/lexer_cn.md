# 什么是SQL解析器
SQL解析器对输入的SQL文本进行切分, 检查是否符合一定的语法规则, 并生成抽象语法树(AST), 供后续优化器进行优化, 它具有词法分析, 语法分析, 语义分析等流程。

图: 数据库体系结构

## 词法分析(lexical analysis)
词法分析按单个字符读取SQL文本, 按照一定词法规则识别各类单词, 将输入文本转换为token序列, 进行词法分析的程序或者函数叫作词法分析器(lexical analyzer, 简称lexer)或者扫描器(scanner).

例: select col1, col2 from tab1 where id = 1+2

经过词法分析, 上面这条SQL语句会切分成select-关键字(keyword), col1-标识符(identifier), col2-标识符(identifier), from-关键字(keyword), tab1-标识符(identifier), where-关键字(keyword), id-标识符(identifier), =-比较运算符(comparisonOperator), 1-数字字面量(numberLiteral), +-算数运算符(arithmeticOperator), 2-数字字面量(numberLiteral)

这里的select, col1, col2, from等被切分出来的输入文本的子串称为词素(lexeme), 关键字, 标识符, 运算符, 字面量等称为token类型, 词素与类型合在一起组成token, 词法分析阶段主要就是对输入的字符流进行适当地切分, 识别token类型, 最后输出token列表, 供后续语法分析使用. 词法分析的难点在于适当地对输入字符进行分隔, 还要考虑一些特殊情况, 有时需要保存当前状态, 必要时进行回溯等情况.

## 语法分析(syntax analysis)
语法分析有时也叫parsing, 它根据上一步生成的token序列检查SQL语句是否符合该数据库支持的语法结构, 它是上下文无关的检查, 进行语法分析的程序或函数叫作语法分析器(syntax analyzer)或者解析器(parser).

例: sel col1 from tab1

在该语句的第一个token为sel, 在仅考虑DML语句时, 第一个token必须为关键字且必须是select, insert, update, delete中的一种, 而sel甚至都不是关键字, 因此会在词法分析阶段被识别为标识符, 在语法分析阶段因不符合SQL语法而报错.

## 语义分析(semantic analysis)
语义分析主要是进行上下文有关的检查, 例如字段类型是否匹配, 变量是否已定义, 表是否存在等, 进行语义分析的程序或函数叫作语义分析器(semantic analyzer).

例: select undefined_func();

该SQL语句符合函数调用的语法规则, 因此可以通过语法分析阶段, 但在进行语义分析时会发现函数undefined_func未定义, 因此会报错.

本文重点在词法分析上, 探讨词法分析器的实现原理.

# 基本原理
## 语句限定
由于SQL语句可以写得非常复杂, 为了便于进行讲解, 本篇将对待解析的SQL语句进行如下限定:
- 仅支持简单select语句, 且句式为select col_name [[as] column_alias][, col_name [[as] column_alias]...] from tab_name where col_name operator expression [and col_name operator expression], 其中中括号的内容为可选项, operator代表操作符, expression代表表达式, 例如: select col1 as alias1, col2 from tab_name where id = 1;
- 关键字仅支持小写英文字符
- 不允许使用反引号(`)转义关键字以当做标识符
- 标识符仅支持小写英文字符, 数字, 下划线
- 字符串字面量仅支持小写英文字符, 数字, 下划线, 且必须用单引号(')包裹
- 数字字面量仅支持整数
- 支持的比较运算符有7种: >, >=, <, <=, = , !=, <>
- 支持表达式, 例如: 'abc', 123, 123 + 456
- 不支持函数, 变量, if else, case when, 子查询, join, 注释, group by, order by, limit, offset等复杂用法
- 多个where条件仅支持and, 不支持or和in
- 当使用空白字符时允许多个相同或不同的空白符连续出现

## 解析流程
词法解析可以用如下步骤描述:
1. 定义词法规范
2. 记R为可以匹配所有词素的正则表达式的集合, 请注意这里的`所有词素`包括不符合词法规范的词素, 其中不符合词法规范的集合可称为error集合, 即R = 关键字+标识符+运算符+字面量....+error = R1+R2+R3...+Rj+...+Rm
3. 记x1x2x3...xn为输入字符流, 记L(R)为符合词法规范的正则表达式的集合, 令1<=i<=n, 检查x1x2x3...xi是否属于L(R)
4. 如果是, 则说明x1x2x3...xi属于L(Rj),Rj属于R
5. 删除x1x2x3...xi, 跳转到步骤3继续匹配

实现方法如下:

图: 词法解析实现方法

## 词法规范(Lexical Specification)
词法规范描述本语言支持的关键字, 标识符, 运算符等token相关的规范, 是词法解析的基础, 根据SQL语法规定以及上面约定的语句限定, 本篇涉及的token类型如下:
- 关键字(keyword): select, as, from, where, and
- 标识符(identifier): 一个或多个小写英文字符或数字或下划线组成, 且不允许由纯数字组成
- 比较运算符(comparisonOperator): >, >=, <, <=, = , !=, <>
- 算数运算符(arithmeticOperator): +, -, *, /, %
- 数字字面量(numberLiteral): 有纯数字组成, 例如: 123等
- 字符串字面量(stringLiteral): 由单引号包裹的小写英文字符或数字或下划线组成, 例如: 'abc', 'aaa_123_bbb'等
- 分隔符(separator): 逗号(,), 分号(;), 左括号((), 右括号()), 单引号(')
- 空白符(whitespace): 空格( ), 回车符(\r), 换行符(\n), 制表符(\t)

具体实现上, token类型会有进一步的细分, 以便后续的语法分析中使用, 例如关键字会进一步细分成selectKeyword, asKeyword, fromKeyword等, 分隔符会细分成LeftBrace, RightBrace等

## 正则表达式(Regular Expression)
确定了词法规范后, 需要将每种token类型用正则表达式来表示, 以便在输入字符流中进行辨识, 根据SQL语法规定以及上面约定的语句限定, 本篇涉及的正则表达式如下:
- 关键字: select|as|from|where|and
- 标识符: [0-9]\*[a-z_][a-z0-9_]*
- 比较运算符: >|>=|<|<=|=|!=|<>
- 算数运算符: +|-|*|/|%
- 数字字面量: [0-9]+
- 字符串字面量: '[a-z0-9_]*'
- 分隔符: ,|;|(|)|'
- 空白符: 空格|\r|\n|\t

其中|表示或, +表示一个或多个, *表示0个或多个

#### 歧义性
针对输入字符流, 如果存在多种切分方法, 则说明存在歧义性(或者叫二义性), 主要存在2种可能性:
- 可按不同长度进行切分

  假设x1x2...xj与x1x2...xk都属于L(R), 且j<k, 则应以更长的k为切分点

  例: select col1 from tab1 where id >= 1

  这里`>=`可以有两种切分方法, 一种是`>`和`=`2个比较运算符, 另一种是`>=`这1个比较运算符, 由与`>=`比`>`更长, 因此使用后一种切分方法

- 切分长度相同但可以识别为不同的类型

  假设x1x2...xi属于L(Rj), 且x1x2...xi属于L(Rk), 即同一个字符串可以识别为不同的类型, 这时按照优先级来进行识别

  例: select col1 from tab1

  这里的select既符合关键字的正则表达式, 也符合标识符的正则表达式, 应识别为哪种类型呢?

  为了解决这个问题, 需要对token类型定义优先级, 通常来讲关键字的优先级最高, error类型的优先级最低, 因此这里应识别为关键字.

  与一些编程语言不同(如C语言), SQL允许标识符以数字开头, 但不允许由纯数字组成标识符, 因为一旦允许纯数字组成标识符, 将难以与数字字面量进行区分, 而标识符与字面量又很难定义哪个优先级更高, 因此直接从词法规范层面规定标识符不允许纯数字组成, 一旦遇到纯数字一律识别为数字字面量, 以避免歧义性带来的麻烦, 标识符不允许加号(+), 小于号(<), 括号(()等其他token的字符也是同理.

## 非确定有限自动机(Non-deterministic Finite Automata, 简称NFA)
自动机也叫自动状态机, 根据读取的输入字符在不同的状态间进行转换, 有限自动机表示该自动机的输入字符与状态是有限的(或者叫有穷的), 非确定有限自动机表示该有限自动机在某一个状态时, 对同一个输入字符可能会存在不同的状态转换, 即在同一个输入字符下, 下一个状态是不确定的. 所有的正则表达式都可以用NFA来等价表示, 可用NFA来实现正则表达式引擎.

状态定义如下:
- S为初始状态
- A, B, C...等为中间状态, 1<=i<=n, Ai, Bi, Ci
- F为终止状态
- ε代表空跳转, 即本次状态转换不消耗输入字符, ε-move表示从某个状态开始进行若干次ε转换, ε-closure(A)表示从A状态出发, 经过ε-move后所能到达的所有的状态的集合

从初始状态开始, 通过读取字符在不同状态之间进行转换, 当读取完所有字符时, 如果存在一种可能性, 可到达终止状态, 即为解析成功, 否则解析失败.



#### 识别select关键字的NFA

图: NFA-select

- 假设输入字符串为select, 从S状态开始读取字母s以后转换为A1状态, 读取字母e以后转换为A2状态, 以此类推一直到读取字母t后转换为A6状态, 此时所有字符已经读取完毕, A6状态下可以通过ε-move到达F状态(即ε-closure(A6)包含终止状态), 因此本次解析成功.
- 假设输入字符串为selct, 经过依次读取字母s, e, l后进入A3状态, 此时读取字母c, 无法跳转到下一个状态, 而A3即不是终止状态也无法通过ε-move到达终止状态(即ε-closure(A3)不包含终止状态), 因此本次解析失败.

#### 识别标识符的NFA

图: NFA-identifier

- 假设输入字符串为1a2b, 从S状态开始读取数字1以后依旧为S状态, 读取字母a以后转换为A状态, 读取数字2以后依旧为A状态, 读取字母b以后依旧是A状态, 此时所有字符串读取完毕, 通过ε-move可到达F状态, 解析成功.
- 假设输入字符串为123, 从S状态开始读取数字1以后依旧为S状态, 读取数字2以后依旧为S状态, 读取数字3以后依旧为S状态, 此时所有字符读取完毕, 状态依旧是S, S状态既不是终止状态也无法通过ε-move到达终止状态, 解析失败.

#### 识别所有关键字的NFA

图: NFA-keyword

#### 识别select关键字与标识符的NFA

图: NFA-select_identifier

F1代表select关键字, F2代表标识符, 假设输入字符串为selct, 由于关键字具有更高优先级, 因此先走上面的分支, 走之前要先保存当前状态(即S状态)与当前读入情况(当前未读入任何字符), 然后依次读取字母s, e, l后进入A3状态, 读取字母c, 此时无法跳转到下一个状态, 因此需要`回溯`到S状态, 通过ε-move到达B1状态, 重新读取字母s后转换为B2状态, 然后依次读取字母e, l, c, t后依旧为B2状态, 最后通过ε-move到达终止状态F2, 本次解析成功且token类型为标识符. NFA由于对同一个输入存在不同的转换, 因此需要记录分叉点, 必要时`回溯`至分叉点, 再尝试另一条路径, 当一个SQL中存在很多的标识符时, 经常需要回溯, 从而降低性能.

#### 识别as关键字与标识符的NFA

图: NFA-as_identifier

通过不断扩充, 可最终实现识别所有token的NFA.

## 确定有限自动机(Deterministic Finite Automata, 简称DFA)
DFA也是一种有限自动机, 可与NFA进行等价变换, 与NFA的区别如下:
- 在状态与输入字符确定的情况下, 下一个状态是唯一的, 即针对同一个输入字符, 仅存在一条转换路径
- 不存在ε-move

DFA不需要保存临时状态, 不需要`回溯`, 因此性能较高, 但是他的状态数量可能会非常多, 假设某个NFA有N个状态, 那理论上DFA最大可能有2的N次方个状态, 因此会占用较多内存空间.

NFA到DFA的转换可通过子集构造法实现, 步骤如下:
- ε-closure(S), 即从S状态出发, 经过ε-move后所能到达的状态的集合, 记为Set0
- Set0里的各个状态可接受的输入字符的集合记为∑, 假设字符x属于∑, 从Set0里的各个状态出发, 获取读取字符x以后所有能达到的状态的集合记为move(Set0, x), move(Set0, x)的ε-closure记为Set1, 即Set1 = ε-closure(move(Set0, x))
- 针对每个xi, xi属于∑, 获取ε-closure(move(Set0, xi))
- 针对每个新产生的Set, 递归地获取所有ε-closure(move(Set, xi))即可获得等价的DFA

以识别as关键字与标识符的NFA举例, 根据定义可得Set0包含S, B1, 针对输入字符a, Set1包含A1, B2, F2, 如图所示:

图: NFA-DFA_1

重复上述步骤以后, 可得下图:

图: NFA-DFA_2

上图中可看到Set2里的F2被打叉了, 因为一个Set里有多个终止状态说明发生了歧义性, 应根据优先级保留唯一的终止状态, 这里保留了代表关键字的F1. 当把图里的每个Set作为一个整体当作一个状态时, 即获得了等价的DFA, 可以看到该DFA里即没有ε-move, 也没有针对同一个输入字符进行分叉的情况, 特别地, Set4里没有包含终止状态, 说明当输入字符串仅包含数字时, 即不能识别为as关键字也不能识别为标识符, 解析会失败.


# demo
词法分析器可通过Flex等工具来实现, 只要定义好正则表达式, Flex可自动生成对应的词法分析器, 即它是生成程序的程序. NFA到DFA的转换是这类工具的核心功能, 不过从业界来看, 大多数语言的编译器选择了纯手工的方式来打造词法分析器, 以更多地进行针对性的优化来提高性能. 不同的编译器或解释器选择了不同的有限自动机来实现词法分析, 这里介绍为本篇文章而编写的SQL解析器的用法.

## 编译安装
```shell
go version
```
go需要1.16以上版本
```shell
git clone https://github.com/romberli/sql-parser-go.git
cd sql-parser-go
go build -o parser main.go
```

## 使用方法
#### 使用NFA进行解析
```
./parser nfa --sql="select col1 as c1, col2 from t01 where id <= 100 and col1 = 'abc'"
```
输出如下:
```text
{tokenType: SelectKeyword, lexeme: select}
{tokenType: Identifier, lexeme: col1}
{tokenType: AsKeyword, lexeme: as}
{tokenType: Identifier, lexeme: c1}
{tokenType: Comma, lexeme: ,}
{tokenType: Identifier, lexeme: col2}
{tokenType: FromKeyword, lexeme: from}
{tokenType: Identifier, lexeme: t01}
{tokenType: WhereKeyword, lexeme: where}
{tokenType: Identifier, lexeme: id}
{tokenType: LessOrEqual, lexeme: <=}
{tokenType: NumberLiteral, lexeme: 100}
{tokenType: AndKeyword, lexeme: and}
{tokenType: Identifier, lexeme: col1}
{tokenType: Equal, lexeme: =}
{tokenType: StringLiteral, lexeme: 'abc'}
```
#### 使用DFA进行解析
```
./parser dfa --sql="select col1 as c1, col2 from t01 where id <= 100 and col1 = 'abc'"
```
输出如下:
```text
{tokenType: SelectKeyword, lexeme: select}
{tokenType: Identifier, lexeme: col1}
{tokenType: AsKeyword, lexeme: as}
{tokenType: Identifier, lexeme: c1}
{tokenType: Comma, lexeme: ,}
{tokenType: Identifier, lexeme: col2}
{tokenType: FromKeyword, lexeme: from}
{tokenType: Identifier, lexeme: t01}
{tokenType: WhereKeyword, lexeme: where}
{tokenType: Identifier, lexeme: id}
{tokenType: LessOrEqual, lexeme: <=}
{tokenType: NumberLiteral, lexeme: 100}
{tokenType: AndKeyword, lexeme: and}
{tokenType: Identifier, lexeme: col1}
{tokenType: Equal, lexeme: =}
{tokenType: StringLiteral, lexeme: 'abc'}
```

# 扩展阅读

- 龙书: Compilers: Principles, Techniques, & Tools
- 斯坦福公开课: CS143 Compilers
- 极客时间: 编译原理之美, 编译原理实践课
