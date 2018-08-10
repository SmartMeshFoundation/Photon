import { Component, OnInit, Input, Output, EventEmitter } from '@angular/core';
import { Observable } from 'rxjs/Observable';
import { BehaviorSubject } from 'rxjs/BehaviorSubject';

import { SmartRaidenConfig } from '../../services/smartraiden.config';
import { SmartRaidenService } from '../../services/smartraiden.service';
import { Event, EventsParam } from '../../models/event';


@Component({
    selector: 'app-event-list',
    templateUrl: './event-list.component.html',
    styleUrls: ['./event-list.component.css']
})
export class EventListComponent implements OnInit {

    @Input() eventsParam: EventsParam;
    @Output() activity: EventEmitter<void> = new EventEmitter<void>();

    public events$: Observable<Event[]>;
    public loading = true;

    constructor(
        private smartraidenConfig: SmartRaidenConfig,
        private smartraidenService: SmartRaidenService
    ) { }

    ngOnInit() {
        let fromBlock: number = this.smartraidenConfig.config.block_start;
        let first = true;
        const data_excl = ['event_type', 'block_number', 'timestamp'];
        const firerSub: BehaviorSubject<void> = new BehaviorSubject(null);
        this.events$ = firerSub
            .do(() => this.loading = true)
            .switchMap(() => this.smartraidenService.getBlockNumber())
            .switchMap((latestBlock) => {
                if (fromBlock > latestBlock) {
                    return Observable.of([]);
                }
                const obs = this.smartraidenService.getEvents(
                    this.eventsParam,
                    fromBlock,
                    latestBlock);
                fromBlock = latestBlock + 1;
                return obs;
            })
            .map((events) => events.map((event) =>
                <Event>Object.assign(event,
                    {
                        data: JSON.stringify(Object.keys(event)
                            .filter((k) => data_excl.indexOf(k) < 0)
                            .reduce((o, k) => (o[k] = event[k], o), {}))
                    }
                ))
            )
            .do((newEvents) => {
                if (newEvents.length > 0 && !first) {
                    this.activity.emit();
                }
                first = false;
            })
            .scan((oldEvents, newEvents) => {
                // this scan/reducer agregates new events (since next_block) with old ones,
                // updating next_block if needed. If no next_block previously,
                // it means it fetched all events, so use only newEvents
                newEvents.reverse(); // most recent first
                const events = !first ? [...newEvents, ...oldEvents] :
                    newEvents.length > 0 ? newEvents : oldEvents;
                for (const event of events) {
                    if (event.block_number > 0 && !event.timestamp) {
                        this.smartraidenService.blocknumberToDate(event.block_number)
                            .subscribe((date) => event.timestamp = date);
                    }
                }
                return events;
            }, [])
            .do(() => setTimeout(() => firerSub.next(null), this.smartraidenConfig.config.poll_interval))
            .do(() => this.loading = false);
    }
}
