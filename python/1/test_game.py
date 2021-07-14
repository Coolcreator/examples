import sys
import unittest
from contextlib import contextmanager
from game import TicTacToe
from io import StringIO
from unittest.mock import patch

@contextmanager
def captured_output():
    new_out, new_err = StringIO(), StringIO()
    old_out, old_err = sys.stdout, sys.stderr
    try:
        sys.stdout, sys.stderr = new_out, new_err
        yield sys.stdout, sys.stderr
    finally:
        sys.stdout, sys.stderr = old_out, old_err

class TestGame(unittest.TestCase):
    def setUp(self):
        self.game = TicTacToe()

    def test_winner(self):
        self.game.board = ['X', 'O', 'X', 'O', 'X', 'O', 'X', ' ', ' ']
        self.assertTrue(self.game.check_winner('X'))
        self.game.board = [' ', ' ', 'X', 'O', 'X', 'O', 'X', 'O', 'X']
        self.assertTrue(self.game.check_winner('X'))
        self.game.board = ['X', ' ', 'X', 'O', 'X', ' ', 'X', 'O', 'O']
        self.assertTrue(self.game.check_winner('X'))
        self.game.board = [' ', 'X', 'O', 'X', 'X', 'O', ' ', ' ', 'O']
        self.assertTrue(self.game.check_winner('O'))
        self.game.board = ['O', 'O', 'O', ' ', 'X', ' ', 'X', ' ', 'X']
        self.assertTrue(self.game.check_winner('O'))
        self.game.board = ['X', ' ', 'O', 'X', 'O', ' ', 'O', 'X', ' ']
        self.assertTrue(self.game.check_winner('O'))

        self.game.board = ['X', 'O', 'X', 'O', 'X', 'O', 'X', ' ', ' ']
        self.assertFalse(self.game.check_winner('O'))
        self.game.board = [' ', ' ', 'X', 'O', 'X', 'O', 'X', 'O', 'X']
        self.assertFalse(self.game.check_winner('O'))
        self.game.board = ['X', ' ', 'X', 'O', 'X', ' ', 'X', 'O', 'O']
        self.assertFalse(self.game.check_winner('O'))
        self.game.board = [' ', 'X', 'O', 'X', 'X', 'O', ' ', ' ', 'O']
        self.assertFalse(self.game.check_winner('X'))
        self.game.board = ['O', 'O', 'O', ' ', 'X', ' ', 'X', ' ', 'X']
        self.assertFalse(self.game.check_winner('X'))
        self.game.board = ['X', ' ', 'O', 'X', 'O', ' ', 'O', 'X', ' ']
        self.assertFalse(self.game.check_winner('X'))

    def test_tie(self):
        self.game.board = ['O', 'X', 'O', 'X', 'O', 'X', 'X', 'O', 'X']
        self.assertTrue(self.game.check_tie())
        self.game.board = ['O', 'X', 'O', 'X', 'O', 'O', 'X', 'O', 'X']
        self.assertTrue(self.game.check_tie())
        self.game.board = ['X', 'O', 'X', 'O', 'X', 'O', 'X', ' ', ' ']
        self.assertFalse(self.game.check_tie())
        self.game.board = ['X', ' ', 'O', 'X', 'O', ' ', 'O', 'X', ' ']
        self.assertFalse(self.game.check_tie())

    def test_reset(self):
        self.game.board = ['X', 'O', 'X', 'O', 'X', 'O', 'X', ' ', ' ']
        self.assertEqual(self.game.reset_game(), 9 * [' '])
        self.game.board = [' ', ' ', 'X', 'O', 'X', 'O', 'X', 'O', 'X']
        self.assertEqual(self.game.reset_game(), 9 * [' '])
        self.game.board = ['O', 'O', 'X', 'X', 'X', 'O', ' ', ' ', ' ']
        self.assertEqual(self.game.reset_game(), 9 * [' '])

    def test_update(self):
        self.game.reset_game()
        self.assertEqual(self.game.update_board(0, 'X'), 'X')
        self.assertEqual(self.game.update_board(1, 'X'), 'X')
        self.assertEqual(self.game.update_board(2, 'O'), 'O')
        self.assertEqual(self.game.update_board(3, 'O'), 'O')
        self.assertEqual(self.game.update_board(4, 'X'), 'X')
        self.assertEqual(self.game.update_board(5, 'X'), 'X')
        self.assertEqual(self.game.update_board(6, 'O'), 'O')
        self.assertEqual(self.game.update_board(7, 'O'), 'O')
        self.assertEqual(self.game.update_board(8, 'X'), 'X')
    
    def test_validate(self):
        self.game.reset_game()
        self.game.board[4] = 'X'
        self.game.board[5] = 'X'
        self.game.board[6] = 'X'
        self.assertTrue(self.game.validate('1'))
        self.assertEqual(self.game.output_message, '')
        self.assertTrue(self.game.validate('2'))
        self.assertEqual(self.game.output_message, '')
        self.assertTrue(self.game.validate('3'))
        self.assertEqual(self.game.output_message, '')
        self.assertFalse(self.game.validate('11'))
        self.assertEqual(self.game.output_message, '(X choise) Please enter a number from the range 1-9: ')
        self.assertFalse(self.game.validate('13'))
        self.assertEqual(self.game.output_message, '(X choise) Please enter a number from the range 1-9: ')
        self.assertFalse(self.game.validate('15'))
        self.assertEqual(self.game.output_message, '(X choise) Please enter a number from the range 1-9: ')
        self.assertFalse(self.game.validate('5'))
        self.assertEqual(self.game.output_message, '(X choise) The 5 cell is already filled. Choose another cell: ')
        self.assertFalse(self.game.validate('6'))
        self.assertEqual(self.game.output_message, '(X choise) The 6 cell is already filled. Choose another cell: ')
        self.assertFalse(self.game.validate('7'))
        self.assertEqual(self.game.output_message, '(X choise) The 7 cell is already filled. Choose another cell: ')


        
    
    def test_show(self):
        self.game.board = ['X', 'O', 'X', 'O', 'X', 'O', 'X', ' ', ' ']
        board = [
            ' X |   |   ',
            '-----------',
            ' O | X | O ',
            '-----------',
            ' X | O | X ']
        
        with captured_output() as (out, _):
            self.game.show_board()
        rows = out.getvalue().splitlines()

        for i in range(len(board)):
            self.assertEqual(rows[i], board[i])


if __name__ == '__main__':
    unittest.main()
