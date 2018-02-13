import click

from .chain import main as chain_main
from .payment import main as payment_main


commands = click.Group('commands', help="Plasma")
commands.add_command(payment_main, "payment")
commands.add_command(chain_main, "chain")

if __name__ == "__main__":
    commands()
