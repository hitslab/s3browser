// Cynhyrchwyd y ffeil hon yn awtomatig. PEIDIWCH Â MODIWL
// This file is automatically generated. DO NOT EDIT
import {config} from '../models';
import {s3} from '../models';

export function Connect(arg1:config.S3Settings):Promise<void>;

export function Download(arg1:string):Promise<void>;

export function GetConfig():Promise<config.Config>;

export function List(arg1:string):Promise<s3.List>;

export function UploadFile(arg1:string):Promise<void>;

export function UploadFolder(arg1:string):Promise<void>;
