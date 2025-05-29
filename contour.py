import numpy as np

# n개의 자연수 중 2개를 순서 없이 뽑는 모든 경우의 수를 반환
def two_comb(n: int) -> list[list[int, int]]:
    comb = []
    for i in range(n):
        for j in range(i+1, n):
            comb.append([i, j])
    return comb

# 주어진 부분영역의 경계면을 리스트로 반환
def all_faces(shape: list) -> list[list[list[int]]]:
    num = np.prod(shape)
    faces = []
    for i in range(num):
        base_index = np.array(np.unravel_index(i, shape), dtype=int)
        for axes in two_comb(len(shape)):
            base_axes = np.array(base_index, dtype=int)[axes]
            shape_axes = np.array(shape, dtype=int)[axes]-1
            if any(base_axes == shape_axes):
                continue
            face = np.array([base_index for _ in range(4)], dtype=int)
            off = np.zeros((4, len(shape)), dtype=int)
            off[:, axes] = np.array([[0, 0], [0, 1], [1, 1], [1, 0]])
            face += off
            faces.append([list(face[i]) for i in range(4)])
    return faces

# 등위선과 등위면을 그리는 함수를 구현
def contour(mesh: np.array, lv: float) -> np.array:
    lines = np.empty((0, mesh.shape[0]-1)) # 선분의 끝점 저장 변수
    for face in all_faces(mesh.shape[1:]):
        rs = [] # 부분영역의 모서리에서 함숫값이 lv인 점의 정보를 저장할 변수
        for i in range(4):            
            fp, fq = mesh[-1, *face[i]], mesh[-1, *face[(i+1)%4]]
            if (fp - lv)*(fq - lv) <= 0:
                t = (lv - fp)/(fq - fp)
                p, q = mesh[:-1, *face[i]], mesh[:-1, *face[(i+1)%4]]
                rs.append((1-t)*p + t*q)
        for i in range(len(rs)):
            for j in range(i+1, len(rs)):
                lines = np.vstack((lines, # 변수의 마지막 행에 선분의
                                   np.array([rs[i], rs[j]]), # 양 끝점을 추가하고
                                   np.full((1, mesh.shape[0]-1), np.nan))) # NaN 삽입
    return lines