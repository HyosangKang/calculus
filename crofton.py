import numpy as np

# 함수 intersections의 구현
def intersections(arry: np.array) -> int:
    num_intrs = 0
    is_black = True
    for i in range(len(arry)):
        if arry[i] == 1 and is_black:
            num_intrs += 1
            is_black = False
        elif arry[i] == 0:
            is_black = True
    return num_intrs


# 함수 crofton_length의 구현
step = 10
def crofton_length(screen) -> float:
    num_intrs = 0
    img = pygame.surfarray.pixels3d(screen)
    height, width = screen.get_size()
    non = np.zeros((height, width))
    for i in range(height):
        for j in range(width):
            if 0 <= i <= tbox_w and 0 <= j <= tbox_h: continue
            non[i, j] = 1 if img[i, j][0] == 255 else 0
    for i in range(0, width, step):
        num_intrs += intersections(non[:, i])
    for j in range(0, height, step):
        num_intrs += intersections(non[j, :])
    return num_intrs * np.pi/2 * 0.1 /2

# 마우스로 그린 곡선의 길이를 출력하는 프로그램 제작
import pygame

pygame.init()
pygame.font.init()
font = pygame.font.SysFont('Comic Sans MS', 30)
screen = pygame.display.set_mode((800, 600))
pygame.draw.rect(screen, (255, 255, 255), (0, 0, tbox_w, tbox_h))

play = True
while play:
    for event in pygame.event.get():
        if event.type == pygame.QUIT:
            play = False
        if event.type == pygame.MOUSEBUTTONDOWN:
            if pygame.mouse.get_pressed()[0]:
                x, y = pygame.mouse.get_pos()
                pygame.draw.circle(screen, (255, 255, 255), (x, y), 3)
            if pygame.mouse.get_pressed()[2]:
                screen.fill((0, 0, 0))
        # 마우스 버튼을 떼면 길이를 출력
        if event.type == pygame.MOUSEBUTTONUP:
            length = crofton_length(screen)
            text = font.render(f'%.1f'%length, True, (0, 0, 0))
            screen.fill((255, 255, 255), (0, 0, tbox_w, tbox_h))
            screen.blit(text, (20, 5))
    if pygame.mouse.get_pressed()[0]:
        new_x, new_y = pygame.mouse.get_pos()
        pygame.draw.line(screen, (255, 255, 255), (x, y), (new_x, new_y), 5)
        x, y = new_x, new_y
    pygame.display.flip()

pygame.quit()

