import os

class TicTacToe:

    def __init__(self):
        self.board = 9 * [' ']
        self.player = 'X'
        self.output_message = ''

    def refresh_screen(self):
        os.system('clear')
        self.show_board()

    def show_board(self):
        print(' %s | %s | %s ' %(self.board[6], self.board[7], self.board[8]))
        print('-----------')
        print(' %s | %s | %s ' %(self.board[3], self.board[4], self.board[5]))
        print('-----------')
        print(' %s | %s | %s ' %(self.board[0], self.board[1], self.board[2]))

    def update_board(self, board_cell, player):
        self.board[board_cell] = player
        return self.board[board_cell]

    def check_winner(self, player):
        for cells in [[6, 4, 2], [8, 4, 0],
                      [6, 7, 8], [3, 4, 5], [0, 1, 2],
                      [6, 3, 0], [7, 4, 1], [8, 5, 2]]:
            result = True
            for cell in cells:
                if self.board[cell] != player:
                    result = False
            if result:
                self.output_message = f'\nThe winner is {self.player}!\n'
                return True
        return False

    def check_tie(self):
        for cell in self.board:
            if cell == ' ':
                return False
        self.output_message = '\nTie game!\n'
        return True



    def validate(self, cell):
        if not cell.isdigit():
            self.output_message = f'({self.player} choise) ' +\
                            f'Must be integer from the range 1-9: '
            return False
        cell = int(cell)
        if cell < 1 or cell > 9:
            self.output_message = f'({self.player} choise) ' +\
                        f'Please enter a number from the range 1-9: '
            return False
        if self.board[cell-1] != ' ':
            self.output_message = f'({self.player} choise) ' +\
                        f'The {cell} cell is already filled. Choose another cell: '
            return False
        return True



    def check_input(self):
        self.output_message = f'\n({self.player} choise) Please enter cell number: '
        while True:
            cell = input(self.output_message)
            ok = self.validate(cell)
            cell = int(cell)
            if ok:
                return cell
            continue

    def start_game(self):
        while True:
            self.refresh_screen()
            cell = self.check_input()
            self.update_board(cell-1, self.player)
            self.refresh_screen()
            if self.check_winner(self.player) or self.check_tie():
                answer = input(self.output_message +
                               'Would you like to play again?\n' +
                               'Type "y" if you want to continue ' +
                               'or any key to stop the game: ')
                if answer == 'y':
                    self.reset_game()
                    continue
                break
            self.player = 'O' if self.player == 'X' else 'X'

    def reset_game(self):
        self.board = 9 * [' ']
        self.player = 'X'
        self.output_message = ''
        return self.board


if __name__ == '__main__':
    game = TicTacToe()
    game.start_game()
