$module('nominalbase', function () {
  var $static = this;
  this.$class('fc0ceb02', 'SomeClass', false, '', function () {
    var $static = this;
    var $instance = this.prototype;
    $static.new = function () {
      var instance = new $static();
      instance.SomeField = $t.fastbox(true, $g.________testlib.basictypes.Boolean);
      return instance;
    };
    this.$typesig = function () {
      return {
      };
    };
  });

  this.$type('ad913192', 'FirstNominal', false, '', function () {
    var $instance = this.prototype;
    var $static = this;
    this.$box = function ($wrapped) {
      var instance = new this();
      instance[BOXED_DATA_PROPERTY] = $wrapped;
      return instance;
    };
    this.$roottype = function () {
      return $g.nominalbase.SomeClass;
    };
    $instance.SomeProp = $t.property(function () {
      var $this = this;
      return $t.fastbox(!$this.$wrapped.SomeField.$wrapped, $g.________testlib.basictypes.Boolean);
    });
    this.$typesig = function () {
      if (this.$cachedtypesig) {
        return this.$cachedtypesig;
      }
      var computed = {
        "SomeProp|3|54ff3ddf": true,
      };
      return this.$cachedtypesig = computed;
    };
  });

  this.$type('eb757a50', 'SecondNominal', false, '', function () {
    var $instance = this.prototype;
    var $static = this;
    this.$box = function ($wrapped) {
      var instance = new this();
      instance[BOXED_DATA_PROPERTY] = $wrapped;
      return instance;
    };
    this.$roottype = function () {
      return $g.nominalbase.SomeClass;
    };
    $instance.GetValue = function () {
      var $this = this;
      return $t.fastbox(!$t.box($this, $g.nominalbase.FirstNominal).SomeProp().$wrapped, $g.________testlib.basictypes.Boolean);
    };
    this.$typesig = function () {
      if (this.$cachedtypesig) {
        return this.$cachedtypesig;
      }
      var computed = {
        "GetValue|2|fd8bc7c9<54ff3ddf>": true,
      };
      return this.$cachedtypesig = computed;
    };
  });

  $static.TEST = function () {
    var sc;
    var sn;
    sc = $g.nominalbase.SomeClass.new();
    sn = $t.box($t.fastbox(sc, $g.nominalbase.FirstNominal), $g.nominalbase.SecondNominal);
    return sn.GetValue();
  };
});
