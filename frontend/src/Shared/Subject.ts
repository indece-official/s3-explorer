export class Subject<T>
{
    private readonly _subscriptions: {[key: number]: ( t: T ) => any};
    private readonly _subscriptionRefs: WeakMap<Object, number>;
    private _nextID: number;


    constructor ( )
    {
        this._subscriptions = {};
        this._subscriptionRefs = new WeakMap();
        this._nextID = 1;
    }


    public subscribe ( owner: Object, clb: ( t: T ) => any ): void
    {
        const id = this._nextID;

        this._subscriptionRefs.set(owner, id);
        this._subscriptions[id] = clb;

        this._nextID++;
    }


    public unsubscribe ( owner: Object ): void
    {
        const id = this._subscriptionRefs.get(owner);
        if ( ! id )
        {
            return;
        }

        if ( this._subscriptions[id] )
        {
            delete this._subscriptions[id];
        }

        this._subscriptionRefs.delete(owner);
    }


    public next ( t: T ): void
    {
        for ( const id in this._subscriptions )
        {
            this._subscriptions[id](t);
        }
    }
}