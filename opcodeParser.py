import json

with open('opcodes.json', 'r') as f:
    data = json.load(f)
    # print(json.dumps(data))
    def op_LD():
        instructions = {}
        for opcode, details in data["unprefixed"].items():
            if details["mnemonic"] == "LD":
                instructions[opcode] = [
                    {"name":operand["name"], "immediate": str(operand["immediate"]) }
                for operand in details.get("operands", [])]
    #             operands = [operand["name"] for operand in details.get("operands", [])]

    #             instructions[opcode] = operands
        return json.dumps(instructions)

    def op_ADD():
            instructions = {}
            for opcode, details in data["unprefixed"].items():
                if details["mnemonic"] == "ADD":
                    instructions[opcode] = [
                        {"name":operand["name"], "immediate": str(operand["immediate"]), "flags": details.get("flags", {}) }
                    for operand in details.get("operands", [])]
            return json.dumps(instructions)

    def op_INC():
        instructions = {}
        for opcode, details in data["unprefixed"].items():
            if details["mnemonic"] == "INC":

                instructions[opcode] = [
                    {"name":operand["name"], "immediate": str(operand["immediate"]), "flags": details.get("flags", {})  }
                for operand in details.get("operands", [])]
        return json.dumps(instructions)
# print(json.dumps(instructions))

