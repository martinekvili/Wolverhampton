using Microsoft.Build.Framework;
using Microsoft.Build.Utilities;
using System;
using System.Collections.Generic;
using System.IO;
using System.Linq;
using System.Runtime.Serialization;
using System.Text;
using System.Threading.Tasks;
using System.Xml;

namespace XmlLogger
{
    public class XmlLogger : Logger
    {
        private BuildResult result;

        public XmlLogger()
        {
            result = new BuildResult();
        }

        public override void Initialize(IEventSource eventSource)
        {
            eventSource.ErrorRaised += eventSource_ErrorRaised;
            eventSource.WarningRaised += eventSource_WarningRaised;
            eventSource.BuildFinished += eventSource_BuildFinished;
        }

        void eventSource_BuildFinished(object sender, BuildFinishedEventArgs e)
        {
            result.Successful = e.Succeeded;
        }

        void eventSource_WarningRaised(object sender, BuildWarningEventArgs e)
        {
            result.ErrorList.Add(new Error(e));
        }

        void eventSource_ErrorRaised(object sender, BuildErrorEventArgs e)
        {
            result.ErrorList.Add(new Error(e));
        }

        public override void Shutdown()
        {
            var dcs = new DataContractSerializer(typeof(BuildResult));
            var settings = new XmlWriterSettings{ Indent = true };
            using (var writer = XmlWriter.Create(@"buildresult.xml", settings))
            {
                writer.WriteStartDocument();
                dcs.WriteObject(writer, result);
            }
        }
    }
}
