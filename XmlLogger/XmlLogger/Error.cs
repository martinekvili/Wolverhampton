using Microsoft.Build.Framework;
using System;
using System.Collections.Generic;
using System.Linq;
using System.Runtime.Serialization;
using System.Text;
using System.Threading.Tasks;

namespace XmlLogger
{
    [DataContract(Name = "error")]
    internal class Error
    {
        [DataContract(Name = "errortype")]
        public enum ErrorType
        {
            [EnumMember(Value = "error")]
            Error, 
            [EnumMember(Value = "warning")]
            Warning
        }

        [DataMember(Name = "type", Order = 0)]
        public ErrorType Type;

        [DataMember(Name = "filename", Order = 1)]
        public string FileName;
        [DataMember(Name = "linenumber", Order = 2)]
        public int LineNumber;
        [DataMember(Name = "columnnumber", Order = 3)]
        public int ColumnNumber;

        [DataMember(Name = "code", Order = 4)]
        public string Code;
        [DataMember(Name = "message", Order = 5)]
        public string Message;

        public Error() {}

        public Error(BuildErrorEventArgs e)
        {
            Type = ErrorType.Error;

            FileName = e.File;
            LineNumber = e.LineNumber;
            ColumnNumber = e.ColumnNumber;

            Code = e.Code;
            Message = e.Message;
        }

        public Error(BuildWarningEventArgs e)
        {
            Type = ErrorType.Warning;

            FileName = e.File;
            LineNumber = e.LineNumber;
            ColumnNumber = e.ColumnNumber;

            Code = e.Code;
            Message = e.Message;
        }
    }
}
