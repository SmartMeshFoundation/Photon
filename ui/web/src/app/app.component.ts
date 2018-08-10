import { Component, OnInit } from '@angular/core';
import { SmartRaidenService } from './services/smartraiden.service';
import { SharedService } from './services/shared.service';

@Component({
    selector: 'app-root',
    templateUrl: './app.component.html',
    styleUrls: ['./app.component.css'],
})
export class AppComponent implements OnInit {
    public title = 'SmartRaiden';
    public smartraidenAddress;
    public menuCollapsed = false;

    constructor(public sharedService: SharedService,
                public smartraidenService: SmartRaidenService) { }

    ngOnInit() {
        this.smartraidenService.getSmartRaidenAddress()
            .subscribe((address) => this.smartraidenAddress = address);
    }

}
