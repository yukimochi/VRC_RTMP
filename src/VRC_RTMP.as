package {

import flash.display.Sprite;
import flash.display.LoaderInfo;
import flash.display.StageScaleMode;
import flash.display.StageAlign;
import flash.events.NetStatusEvent;
import flash.net.NetConnection;
import flash.net.NetStream;
import flash.media.Video;

[SWF(backgroundColor="0x000000")]
public class VRC_RTMP extends Sprite {

    private var nc:NetConnection;
    private var ns:NetStream;
    private var video:Video;

    function VRC_RTMP() {
        var flashvars = LoaderInfo(this.loaderInfo).parameters;
        setupStage();
        setupNetConnection();
        setupVideo();
        nc.connect(flashvars.addr);
    }

    private function setupStage():void {
        stage.scaleMode = StageScaleMode.NO_SCALE;
        stage.align = StageAlign.TOP_LEFT;
    }

    private function setupNetConnection():void {
        nc = new NetConnection();
        nc.addEventListener(NetStatusEvent.NET_STATUS, onChangeNCStatus);
    }

    private function setupVideo():void {
        video = new Video(1920, 1080);
        addChild(video);
    }

    private function setupNetStream():void {
        ns = new NetStream(nc);
        ns.addEventListener(NetStatusEvent.NET_STATUS, onChangeNCStatus);

        video.attachNetStream(ns);
    }

    private function onChangeNCStatus(e:NetStatusEvent):void {
        const code:String = e.info.code;
        var flashvars = LoaderInfo(this.loaderInfo).parameters;
        if (code === "NetConnection.Connect.Success") {
            setupNetStream();
            ns.play(flashvars.streamkey);
        }
    }
}

}