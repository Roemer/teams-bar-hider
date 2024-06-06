using System;
using System.Windows;

using System.Diagnostics;
using System.Runtime.InteropServices;
using FlaUI.Core.Definitions;
using UIA = FlaUI.UIA3;

namespace TeamsBarHider
{
    /// <summary>
    /// Interaction logic for MainWindow.xaml
    /// </summary>
    public partial class MainWindow : Window
    {
        private IntPtr nativeWindowHandle = IntPtr.Zero;

        public MainWindow()
        {
            InitializeComponent();
        }

        private void OnHideClick(object sender, RoutedEventArgs e)
        {
            var teamsProcess = Process.GetProcessesByName("ms-teams");
            if (teamsProcess.Length > 0)
            {
                var id = teamsProcess[0].Id;
                var app = FlaUI.Core.Application.Attach(id);
                using (var automation = new UIA.UIA3Automation())
                {
                    var win = automation.GetDesktop().FindFirstChild(cf =>
                    cf.ByProcessId(id).And(
                        cf.ByName("Freigabesteuerungsleiste", PropertyConditionFlags.MatchSubstring).Or(
                            cf.ByName("Sharing control bar", PropertyConditionFlags.MatchSubstring))
                        )
                    );
                    if (win != null)
                    {
                        nativeWindowHandle = win.FrameworkAutomationElement.NativeWindowHandle;
                        ShowWindow(nativeWindowHandle, 0);
                    }
                }
            }
        }

        private void OnShowClick(object sender, RoutedEventArgs e)
        {
            if (nativeWindowHandle != IntPtr.Zero)
            {
                ShowWindow(nativeWindowHandle, 5);
            }
        }

        [DllImport("user32.dll", CharSet = CharSet.Auto, ExactSpelling = true, SetLastError = true)]
        static extern bool ShowWindow(IntPtr hWnd, int nCmdShow);
    }
}
