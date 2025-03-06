import json

with open('opcodes.json', 'r') as f:
    data = json.load(f)

def parse_operations():
    instructions = {}

    for opcode, details in data["unprefixed"].items():
        mnemonic = details["mnemonic"]  

        if mnemonic not in instructions:
            instructions[mnemonic] = {}

        instructions[mnemonic][opcode] = [
            {"name": operand["name"], "immediate": str(operand["immediate"])}
            for operand in details.get("operands", [])
        ]

        if "flags" in details:
            instructions[mnemonic][opcode].append({"flags": details["flags"]})

    return json.dumps(instructions)
