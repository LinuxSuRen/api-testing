import * as jspb from 'google-protobuf'



export class TestTask extends jspb.Message {
  getData(): string;
  setData(value: string): TestTask;

  getKind(): string;
  setKind(value: string): TestTask;

  getCasename(): string;
  setCasename(value: string): TestTask;

  getLevel(): string;
  setLevel(value: string): TestTask;

  getEnvMap(): jspb.Map<string, string>;
  clearEnvMap(): TestTask;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): TestTask.AsObject;
  static toObject(includeInstance: boolean, msg: TestTask): TestTask.AsObject;
  static serializeBinaryToWriter(message: TestTask, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): TestTask;
  static deserializeBinaryFromReader(message: TestTask, reader: jspb.BinaryReader): TestTask;
}

export namespace TestTask {
  export type AsObject = {
    data: string,
    kind: string,
    casename: string,
    level: string,
    envMap: Array<[string, string]>,
  }
}

export class HelloReply extends jspb.Message {
  getMessage(): string;
  setMessage(value: string): HelloReply;

  getError(): string;
  setError(value: string): HelloReply;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): HelloReply.AsObject;
  static toObject(includeInstance: boolean, msg: HelloReply): HelloReply.AsObject;
  static serializeBinaryToWriter(message: HelloReply, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): HelloReply;
  static deserializeBinaryFromReader(message: HelloReply, reader: jspb.BinaryReader): HelloReply;
}

export namespace HelloReply {
  export type AsObject = {
    message: string,
    error: string,
  }
}

export class Empty extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Empty.AsObject;
  static toObject(includeInstance: boolean, msg: Empty): Empty.AsObject;
  static serializeBinaryToWriter(message: Empty, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): Empty;
  static deserializeBinaryFromReader(message: Empty, reader: jspb.BinaryReader): Empty;
}

export namespace Empty {
  export type AsObject = {
  }
}

