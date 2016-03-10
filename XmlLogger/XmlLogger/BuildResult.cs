using System;
using System.Collections.Generic;
using System.Linq;
using System.Runtime.Serialization;
using System.Text;
using System.Threading.Tasks;

namespace XmlLogger
{
    [DataContract(Name = "buildresult")]
    internal class BuildResult
    {
        [DataMember(Name = "successful", Order = 0)]
        public bool Successful;

        [DataMember(Name = "errorlist", Order = 1)]
        public List<Error> ErrorList;

        public BuildResult()
        {
            ErrorList = new List<Error>();
        }
    }
}
