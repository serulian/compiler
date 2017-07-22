$module('indexer', function () {
  var $static = this;
  this.$class('272ade03', 'SomeClass', false, '', function () {
    var $static = this;
    var $instance = this.prototype;
    $static.new = function () {
      var instance = new $static();
      instance.result = $t.fastbox(false, $g.________testlib.basictypes.Boolean);
      return instance;
    };
    $instance.$index = function (someParam) {
      var $this = this;
      return $t.fastbox($this.result.$wrapped && !someParam.$wrapped, $g.________testlib.basictypes.Boolean);
    };
    $instance.$setindex = function (index, value) {
      var $this = this;
      $this.result = value;
      return;
    };
    this.$typesig = function () {
      return {
      };
    };
  });

  $static.TEST = function () {
    var sc;
    sc = $g.indexer.SomeClass.new();
    sc.$setindex($t.fastbox(1, $g.________testlib.basictypes.Integer), $t.fastbox(true, $g.________testlib.basictypes.Boolean));
    return sc.$index($t.fastbox(false, $g.________testlib.basictypes.Boolean));
  };
});
