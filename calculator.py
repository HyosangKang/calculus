# 점화식을 이용해 수열의 일반항을 계산
def sqrt(n: int) -> float:
    term = 1
    for _ in range(n):
        term = (term + 2/term)/2
    return term

# 지수함수의 계산
def exp(x: float) -> float:
    for _ in range(30): x /= 2   # x/2^30을 계산
    x += 1                       # x에 1 + x/2^30을 대입
    for _ in range(30): x *= x   # x^(2^30)을 계산
    return x 

# 이분법을 이용해 로그함수를 구현
def log(x: float) -> float:
    if x <= 0:
        return float('nan')
    n = 1
    while not (exp(-n) <= x <= exp(n)): n += 1
    a, b, c = -n, n, 0
    while not (a == c or b == c):
        if exp(c) >= x: b = c
        else:           a = c
        c = (b+a)/2
    return c

# 멱급수를 이용해 코사인함수를 파이썬 코드로 구현
pi = 3.141592653589793

def cos(x):
    while x <= -pi:                # 입력된 값이 -pi보다 커질 때까지 2pi를 더함
        x += 2*pi
    while x >= pi:                 # 입력된 값이 pi보다 작아질 때까지 2pi를 뺌
        x -= 2*pi  
    sum, term, i = 1, 1, 1
    while term != 0:               # 급수의 항이 0과 같아질 때까지 반복
        term *= x*x/(2*i)/(2*i-1)  # 급수의 n번째 항의 크기를 계산
        if i%2 == 0: sum += term   # 인덱스가 짝수이면 더함
        else:        sum -= term   # 인덱스가 홀수이면 뺌
        i += 1
    return sum

# 멱급수를 이용해 사인함수를 파이썬 코드로 구현
def sin(x):
    while x <= -pi:
        x += 2*pi
    while x >= pi:
        x -= 2*pi        
    sum, term, i = x, x, 1
    while term != 0:
        term *= x*x/(2*i+1)/(2*i)
        if i%2 == 0: sum += term
        else:        sum -= term
        i += 1
    return sum

# 기본함수를 하나의 딕셔너리에 정의
basic_func = {    # 딕셔너리 타입으로 정의
    'sqrt': sqrt, # 문자열(string) 키(key) sqrt에는 sqrt 함수를 값(value)으로 할당
    'exp': exp, 
    'log': log, 
    'sin': sin, 
    'cos': cos,
}

# 사칙연산을 함수로 정의
add = lambda x, y: x + y
subs = lambda x, y: x - y
mul = lambda x, y: x * y
div = lambda x, y: x / y

# 사칙연산을 하나의 딕셔너리에 정의
arithmetic = {
    '+': add,
    '-': subs,
    '*': mul,
    '/': div,
}

val = 0.0

while 1:
    expr = input("> ")      # 프롬프트(prompt) 시작을 가리키는 문자(>)를 출력하고
    try:                    # 그 뒤의 문자열을 입력받은 뒤
        val = float(expr)   # 문자열을 실수형 자료로 형 변환한다.
    except ValueError:      # 만약 변환이 실패하면 문자열이 숫자가 아니라는 의미
        if expr == 'exit':  # 입력된 문자열이 'exit'이면 
            break           # 프로그램 종료
        for name, func in basic_func.items():    # 기본함수의 이름 중
            if expr == name:                     # 입력된 문자가 일치하면 
                val = func(val)                  # 함숫값을 계산
        for name, func in arithmetic.items():    # 사칙연산 기호 중
            if expr == name:                     # 입력된 문자가 일치하면
                new = input("> ")                # 다음 입력값을 기다린 뒤
                try:                             # 두 수에 대해
                    val = func(val, float(new))  # 해당 연산을 수행
                except ValueError:  # 만약 두 번째 입력된 값이 숫자 형태가 아니면
                    pass            # 연산을 무시
    print(val)